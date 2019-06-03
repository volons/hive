package platform

import (
	"errors"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/volons/hive/controllers/user"
	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/libs/pubsub"
	"github.com/volons/hive/libs/store"
	"github.com/volons/hive/libs/websocket"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

type platform struct {
	*pubsub.Topic

	parser     messages.Parser
	client     *websocket.Client
	status     Status
	statusLock sync.RWMutex
	url        *url.URL
	reconnect  bool
}

// Platform allows communicating with the API
var Platform = &platform{
	Topic:  pubsub.NewTopic(),
	status: Status{},
	parser: messages.NewParser(func(typ string) interface{} {
		switch typ {
		case "position", "goto":
			return &models.Position{}
		case "battery":
			return &models.Battery{}
		case "rc":
			return &models.Rc{}
		case "fence":
			return &models.Fence{}
		case "webrtc:sdp":
			return &models.SessionDescription{}
		case "webrtc:icecandidate":
			return &models.IceCandidate{}
		case "webrtc:start", "takeoff", "land", "rtl":
			return &struct{}{}
		case "caps":
			return &models.Caps{}
		default:
			return &libs.JSONObject{}
		}
	}),
}

// Run connects to the API and listens for messages
func (p *platform) Run(urlBase string) {
	var err error
	p.url, err = url.Parse(urlBase)
	if err != nil {
		log.Println("Bad url", urlBase)
		return
	}

	p.reconnect = true

	p.statusLock.Lock()
	p.status = Status{}
	p.statusLock.Unlock()

	for p.reconnect {
		p.client = websocket.NewClient(false, p.parser)
		err := p.client.Connect(p.url.String())

		if err == nil {
			log.Println("connected to API")
			p.listen()
		} else {
			log.Println("Could not connect to API", err)
		}

		p.statusLock.Lock()
		p.status.Connected = false
		p.status.ID = ""
		p.statusLock.Unlock()

		p.statusUpdated()
		time.Sleep(1 * time.Second)
	}
}

func (p *platform) setToken(token string) {
	q := p.url.Query()
	q.Set("token", token)
	p.url.RawQuery = q.Encode()
}

func (p *platform) statusUpdated() {
	p.Publish(p.status)
}

func (p *platform) GetStatus() Status {
	p.statusLock.RLock()
	defer p.statusLock.RUnlock()
	return p.status
}

// GetFBLiveURL requests a facebook live stream url from the API
func (p *platform) GetFBLiveURL(fbToken string, text string) (libs.JSONObject, error) {
	if p.client == nil {
		return nil, errors.New("Not connected to API")
	}

	return p.Request(messages.New("cmd:get_live_url", libs.JSONObject{
		"access_token": fbToken,
		"text":         text,
	}), time.Second*50)
}

// OpenLocation sends informations about the location
// to the API and sets it as open
func (p *platform) OpenLocation(data libs.JSONObject) error {
	_, err := p.Request(messages.New("open", data), time.Second*20)

	if err != nil {
		return err
	}

	return nil
}

// ListenToQueue subscribes to the users queue sends updates on returned channel
func (p *platform) QueueSubscribe() error {
	return p.Send(messages.New("queue:subscribe", libs.JSONObject{}))
}

// QueuePick picks a user by id in the queue and sends him through to the cockpit
func (p *platform) QueuePick(id, token string) (*models.User, error) {
	res, err := p.Request(messages.New("queue:pick", libs.JSONObject{
		"id":    id,
		"token": token,
	}), time.Second*20)

	if err != nil {
		return nil, err
	}

	user := store.Users.Get(token)
	if user == nil {
		return nil, errors.New("queue pick failed")
	}

	name, ok := res.GetString("name")
	if ok {
		user.SetName(name)
	}

	return user, nil
}

// QueueNext picks the first user in the queue
func (p *platform) QueueNext(token string) (*models.User, error) {
	res, err := p.Request(messages.New("queue:next", libs.JSONObject{
		"token": token,
	}), time.Second*20)

	if err != nil {
		return nil, err
	}

	user := store.Users.Get(token)
	if user == nil {
		return nil, errors.New("queue pick failed")
	}

	name, ok := res.GetString("name")
	if ok {
		user.SetName(name)
	}

	return user, nil
}

func (p *platform) listen() {
	start := time.Now()
	for {
		select {
		case message := <-p.client.Recv():
			p.onMessage(message)
		case <-p.client.Done():
			log.Println("disconnected from API after", time.Since(start))
			return
		}
	}
}

func (p *platform) onMessage(msg messages.Message) {
	switch msg.Type {
	case "login":
		p.onLogin(msg)
	case "connected":
		p.onUserConnected(msg)
	case "disconnected":
		p.onUserDisconnected(msg)
	case "update:queue":
		p.onQueueUpdate(msg)
	case "fwd":
		p.onForwardedMessage(msg)
	case "error":
		p.onError(msg)
		//case "reply":
		//	p.onReply(msg)
	}
}

func (p *platform) onLogin(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		log.Println("login msg without data")
		return
	}

	id, _ := data.GetString("id")

	token, _ := data.GetString("token")
	p.setToken(token)

	p.statusLock.Lock()
	p.status.Connected = true
	p.status.ID = id
	p.status.Token = token
	p.statusLock.Unlock()

	p.statusUpdated()
}

func (p *platform) onUserConnected(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		return
	}

	token, ok := data.GetString("token")
	if !ok {
		p.DisconnectUser(token, errors.New("need token"))
		return
	}

	usr, err := store.Users.Authenticate(token)
	if err != nil {
		p.DisconnectUser(token, err)
		return
	}

	usrConn := user.NewUserController(newFwdClient(token, p.client))
	go usrConn.Start(usr, nil)
}

func (p *platform) onUserDisconnected(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		return
	}

	token, ok := data.GetString("token")
	if ok {
		fwdCli := GetFwdClient(token)
		if fwdCli != nil {
			fwdCli.Stop()
		}
	}
}

func (p *platform) onForwardedMessage(fwdMsg messages.Message) {
	data := fwdMsg.JSONData()
	if data == nil {
		log.Println("no data to fwd")
		return
	}

	from, ok := data.GetString("from")
	if !ok {
		log.Println("no sender of fwd msg")
		return
	}

	msg, ok := data.GetString("msg")
	if !ok {
		log.Println("no msg to fwd")
		return
	}

	fwdCli := GetFwdClient(from)
	if fwdCli == nil {
		log.Println("fwd client not found for token", from)
		p.DisconnectUser(from, errors.New("not authorized"))
		return
	}

	fwdCli.OnMessage(msg)
}

func (p *platform) onQueueUpdate(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		return
	}

	list, ok := data.GetArray("list")
	if !ok {
		return
	}

	store.Queue.Set(list)
}

func (p *platform) onError(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		return
	}

	message, _ := data.GetString("message")
	action, ok := data.GetString("action")
	if ok && action == "logout" {
		p.statusLock.Lock()
		p.status.Err = errors.New(message)
		p.statusLock.Unlock()
		p.logout()
	}
}

func (p *platform) logout() {
	p.reconnect = false
	p.client.Disconnect()
}

//func (p *platform) onReply(msg messages.Message) {
//	data := msg.JSONData()
//	if data == nil {
//		return
//	}
//
//	id, ok := data.GetString("id")
//	if !ok {
//		return
//	}
//
//	l, ok := p.listeners[id]
//	if !ok {
//		return
//	}
//
//	select {
//	case l.ch <- &msg:
//	case <-l.done:
//	}
//}

//func (p *platform) addListener(msg messages.Message) listener {
//	l := newListener(msg.ID)
//	p.listeners[l.msgID] = l
//	return l
//}

//func (p *platform) removeListener(l listener) {
//	delete(p.listeners, l.msgID)
//	close(l.done)
//}

func (p *platform) reply(msgID string, result *string, err error) {
	data := libs.JSONObject{
		"id": msgID,
	}

	if err != nil {
		data["error"] = err.Error()
	}
	if result != nil {
		data["result"] = *result
	}

	p.Send(messages.New("reply", data))
}

// DisconnectUser sends a disconnect message to disconnect a user
func (p *platform) DisconnectUser(token string, err error) error {
	data := libs.JSONObject{
		"token": token,
	}

	if err != nil {
		data["error"] = err.Error()
	}

	return p.Send(messages.New("disconnect", data))
}

// Send sends a message to the API
func (p *platform) Send(msg messages.Message) error {
	return p.client.Send(msg)
}

// Request sends a message to the API and waits for a response
func (p *platform) Request(msg messages.Message, timeout time.Duration) (libs.JSONObject, error) {

	cb := callback.New()

	p.Send(msg.ToRequest(cb))
	data, err := cb.Timeout(time.Minute).Wait()

	if err != nil {
		return nil, err
	}

	result, _ := data.(map[string]interface{})

	return libs.JSONObject(result), nil

	//listener := p.addListener(msg)
	//defer p.removeListener(listener)

	//if err := p.Send(msg); err != nil {
	//	return nil, err
	//}

	//select {
	//case msg := <-listener.ch:
	//	data := msg.JSONData()

	//	if err, ok := data.GetString("error"); ok && err != "" {
	//		return nil, errors.New(err)
	//	}

	//	if result, ok := data.GetObj("result"); ok {
	//		return result, nil
	//	}

	//	return nil, nil
	//case <-time.After(timeout):
	//	return nil, errors.New("API request timeout")
	//case <-p.client.Done():
	//	return nil, errors.New("disconnected from API")
	//}
}
