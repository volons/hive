package platform

import (
	"errors"
	"log"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/libs/websocket"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

var connected = make(map[string]*fwdClient)

type fwdClient struct {
	token     string
	ws        *websocket.Client
	messages  chan messages.Message
	callbacks callback.Map
	parser    messages.Parser
	done      libs.Done
}

func GetFwdClient(token string) *fwdClient {
	return connected[token]
}

func newFwdClient(token string, ws *websocket.Client) *fwdClient {
	c := &fwdClient{
		token:     token,
		ws:        ws,
		messages:  make(chan messages.Message),
		callbacks: callback.NewMap(),
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
		done: libs.NewDone(),
	}

	connected[token] = c
	return c
}

func (c *fwdClient) OnMessage(data string) error {
	msg, err := c.parser.Parse([]byte(data))
	if err != nil {
		return err
	}

	if msg.Type == "reply" {
		c.onReply(msg)
		return nil
	}

	if msg.IsRequest() {
		msg.Callback().Listen(func(result interface{}, err error) {
			log.Println("reply called", result, err)

			data := libs.JSONObject{
				"id": msg.ID,
			}
			if err != nil {
				data["error"] = err.Error()
			} else {
				data["result"] = result
			}

			c.Send(messages.New("reply", data))
		})
	}

	c.messages <- msg
	return nil
}

func (c *fwdClient) onReply(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		log.Println("received reply without data")
		return
	}

	id, ok := data.GetString("id")
	if !ok {
		log.Println("received reply without id")
		return
	}

	cb := c.callbacks.Get(id)
	if cb == nil {
		log.Println("received reply for timed out or non existent request")
		return
	}

	log.Println("onReply:", data)

	if err, ok := data.GetString("error"); ok {
		cb.Reject(errors.New(err))
		return
	}

	res, _ := data.GetAny("result")
	cb.Resolve(res)
}

// Sends message to be forwarded to user
func (c *fwdClient) Send(msg messages.Message) error {
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	if msg.IsRequest() {
		c.callbacks.Add(msg.ID, msg.Callback())
	}

	fwd := messages.New("fwd", libs.JSONObject{
		"to":  c.token,
		"msg": data,
	})

	return c.ws.Send(fwd)
}

func (c *fwdClient) Recv() <-chan messages.Message {
	return c.messages
}

func (c *fwdClient) Done() <-chan bool {
	return c.done.WaitCh()
}

func (c *fwdClient) Disconnect() {
	Platform.DisconnectUser(c.token, nil)
	c.Stop()
}

func (c *fwdClient) Stop() {
	delete(connected, c.token)
	c.done.Done()
}
