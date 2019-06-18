package admin

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/autopilot"
	"github.com/volons/hive/libs/store"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
	"github.com/volons/hive/platform"
)

type cmd func(messages.Message) (interface{}, error)

// Admin represents an admin user that can manage
// vehicles and missions
type Admin struct {
	ch                messages.Channel
	admin             *models.Admin
	channels          map[string]*messages.Line
	lastSentTelemetry time.Time
}

// New creates a new admin
func New(ch messages.Channel) *Admin {
	c := &Admin{
		ch:       ch,
		channels: make(map[string]*messages.Line),
	}
	return c
}

// Start initializes the admin connection and starts
// listening for messages
func (a *Admin) Start(admin *models.Admin, err error) {
	if err != nil {
		sendErr := a.sendError(err, "logout")
		if sendErr != nil {
			log.Println(sendErr)
		}

		a.ch.Disconnect()
		return
	}

	a.admin = admin

	sendErr := a.ch.Send(messages.New("login", libs.JSONObject{
		"id": a.admin.ID(),
	}))
	if sendErr != nil {
		log.Println(sendErr)
	}

	a.onPlatformStatus(platform.Platform.GetStatus())
	a.onVehicleListChanged()
	a.sendTelemetry()

	a.run()
}

func (a *Admin) run() {
	vehiclesSub := store.Vehicles.Subscription()
	defer store.Vehicles.Unsubscribe(vehiclesSub)

	telemetrySub := time.NewTicker(time.Millisecond * 800)
	defer telemetrySub.Stop()

	usersSub := store.Users.Subscription()
	defer store.Users.Unsubscribe(usersSub)

	queueSub := store.Queue.Subscription()
	defer store.Queue.Unsubscribe(queueSub)

	platformSub := platform.Platform.Subscription()
	defer platform.Platform.Unsubscribe(platformSub)

	for {
		select {
		case msg := <-a.ch.Recv():
			a.onMessage(msg)

		case <-vehiclesSub.Recv():
			a.onVehicleListChanged()
		case <-telemetrySub.C:
			a.sendTelemetry()
		case <-usersSub.Recv():
			a.onUsersChanged()
		case <-queueSub.Recv():
			a.onQueueChanged()
		case data := <-platformSub.Recv():
			status := data.(platform.Status)
			a.onPlatformStatus(status)

		case <-a.admin.Done():
			a.ch.Disconnect()
		case <-a.ch.Done():
			a.closeChannels()
			return
		}
	}
}

func (a *Admin) onMessage(msg messages.Message) {
	switch msg.Type {
	case "location:open":
		a.exec(a.openLocation, msg)
	case "fence:set":
		a.exec(a.setFence, msg)
	case "fence:enable":
		a.exec(a.enableFence, msg)
	case "fence:disable":
		a.exec(a.disableFence, msg)
	case "permissions:set":
		a.exec(a.setPermissions, msg)
	case "queue:subscribe":
		a.exec(a.queueSubscribe, msg)
	case "queue:pick":
		a.exec(a.queuePick, msg)
	case "queue:next":
		a.exec(a.queueNext, msg)
	case "channel:open":
		a.exec(a.openChannel, msg)
	case "channel:close":
		a.exec(a.closeChannel, msg)
	case "channel:send":
		a.exec(a.sendOnChannel, msg)
	default:
		a.exec(a.unknownMessage, msg)
	}
}

func (a *Admin) exec(fn cmd, msg messages.Message) {
	go func() {
		res, err := fn(msg)
		a.reply(msg.ID, res, err)
	}()
}

func (a *Admin) onVehicleListChanged() {
	err := a.ch.Send(messages.New("vehicles", store.Vehicles.JSON()))
	if err != nil {
		log.Println(err)
	}
}

func (a *Admin) sendTelemetry() {
	pos := store.Vehicles.TelemetryJSON(a.lastSentTelemetry)
	if len(pos) > 0 {
		err := a.ch.Send(messages.New("telemetry", pos))
		if err != nil {
			log.Println(err)
		} else {
			a.lastSentTelemetry = time.Now()
		}
	}
}

func (a *Admin) onUsersChanged() {
	err := a.ch.Send(messages.New("users", store.Users.JSON()))
	if err != nil {
		log.Println(err)
	}
}

func (a *Admin) onQueueChanged() {
	err := a.ch.Send(messages.New("queue", store.Queue.JSON()))
	if err != nil {
		log.Println(err)
	}
}

func (a *Admin) onPlatformStatus(status platform.Status) {
	err := a.ch.Send(messages.New("platform:status", status))
	if err != nil {
		log.Println(err)
	}
}

func (a *Admin) queueSubscribe(msg messages.Message) (interface{}, error) {
	return nil, platform.Platform.QueueSubscribe()
}

func (a *Admin) openLocation(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		data = libs.JSONObject{}
	}

	fence := models.GetFence()
	if fence == nil {
		return nil, errors.New("Fence is not set")
	}

	data["fence"] = fence.JSON()
	msg.Data = data

	return nil, platform.Platform.OpenLocation(data)
}

func (a *Admin) setFence(msg messages.Message) (interface{}, error) {
	fence, ok := msg.Data.(*models.FenceData)
	if ok {
		return nil, models.SetFence(*fence)
	}

	log.Println("setFence bad: ", fence)
	return nil, errors.New("bad fence data format")
}

func (a *Admin) enableFence(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("no data")
	}

	vehicleID, ok := data.GetString("vehicleID")
	if !ok {
		return nil, errors.New("need vehicleID")
	}

	ap := autopilot.Get(vehicleID)
	if ap == nil {
		return nil, errors.New("Vehicle does not exist")
	}

	return nil, ap.Push(messages.New("fence:enable", nil))
}

func (a *Admin) disableFence(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("no data")
	}

	vehicleID, ok := data.GetString("vehicleID")
	if !ok {
		return nil, errors.New("need vehicleID")
	}

	ap := autopilot.Get(vehicleID)
	if ap == nil {
		return nil, errors.New("Vehicle does not exist")
	}

	return nil, ap.Push(messages.New("fence:disable", nil))
}

func (a *Admin) setPermissions(msg messages.Message) (interface{}, error) {
	return nil, nil
}

func (a *Admin) openChannel(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("could not open channel: no data in message")
	}

	channelID, ok := data.GetString("channelID")
	if !ok {
		return nil, errors.New("could not open channel: need channelID parameter")
	}

	c := strings.Split(channelID, ":")
	typ, id := c[0], c[1]

	if typ == "vehicle" {
		ap := autopilot.Get(id)
		if ap == nil {
			return nil, fmt.Errorf("unknown vehicle with ID '%s'", id)
		}

		line := messages.NewLine(true)
		ap.ConnectUser(line)
		a.channels[channelID] = line

		go func() {
			for {
				select {
				case msg := <-line.Recv():
					// Wrap message
					msg = messages.New("channel:message", libs.JSONObject{
						"channelID": channelID,
						"message":   msg,
					})
					if a.ch.Send(msg) != nil {
						line.Disconnect()
					}
				case <-line.Done():
					delete(a.channels, channelID)
				}
			}
		}()

		log.Printf("Opened channel to vehicle '%v'\n", id)
	} else {
		return nil, fmt.Errorf("unknown channel type '%s'", typ)
	}

	return "OK", nil
}

func (a *Admin) sendOnChannel(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("could not send on channel: no data in message")
	}

	channelID, ok := data.GetString("channelID")
	if !ok {
		return nil, errors.New("could not send on channel: need channelID parameter")
	}

	message, ok := data.GetObj("message")
	if !ok {
		return nil, errors.New("could not send on channel: need message parameter")
	}

	t, ok := message.GetString("type")
	if !ok {
		return nil, errors.New("could not send on channel: need message.type parameter")
	}

	d, _ := message.GetObj("data")

	m := messages.New(t, d)
	line := a.channels[channelID]

	log.Printf("Sending message on channel: %v", channelID)
	return nil, line.Send(m)
}

func (a *Admin) closeChannel(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("could not close channel: no data in message")
	}

	channelID, ok := data.GetString("channelID")
	if !ok {
		return nil, errors.New("could not close channel: need channelID parameter")
	}

	line := a.channels[channelID]
	if line == nil {
		return map[string]bool{"wasOpen": false}, nil
	}

	line.Disconnect()

	return map[string]bool{"wasOpen": true}, nil
}

func (a *Admin) closeChannels() {
	for _, line := range a.channels {
		line.Disconnect()
	}
}

func (a *Admin) queuePick(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("need userID and vehicleID")
	}

	userID, ok := data.GetString("userID")
	if !ok {
		return nil, errors.New("need userID")
	}

	vehicleID, ok := data.GetString("vehicleID")
	if !ok {
		return nil, errors.New("need vehicleID")
	}

	token := store.Users.GenerateToken(vehicleID)

	return platform.Platform.QueuePick(userID, token)
}

func (a *Admin) queueNext(msg messages.Message) (interface{}, error) {
	data := msg.JSONData()
	if data == nil {
		return nil, errors.New("need vehicleID")
	}

	vehicleID, ok := data.GetString("vehicleID")
	if !ok {
		return nil, errors.New("need vehicleID")
	}

	token := store.Users.GenerateToken(vehicleID)

	return platform.Platform.QueueNext(token)
}

func (a *Admin) unknownMessage(msg messages.Message) (interface{}, error) {
	return nil, fmt.Errorf("Unknown message of type %s", msg.Type)
}

func (a *Admin) reply(msgID string, result interface{}, err error) {
	data := libs.JSONObject{
		"id": msgID,
	}

	if err != nil {
		data["error"] = err.Error()
	}
	if result != nil {
		data["result"] = result
	}

	sendErr := a.ch.Send(messages.New("reply", data))
	if sendErr != nil {
		log.Println(sendErr)
	}
}

func (a *Admin) sendError(err error, action string) error {
	return a.ch.Send(messages.New("error", libs.JSONObject{
		"message": err.Error(),
		"action":  action,
	}))
}
