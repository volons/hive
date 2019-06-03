package sdk

import (
	"errors"
	"log"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/messages"
)

// sdkChannel implements messages.Channel for communication
// between gate and external nodes
type sdkChannel struct {
	i         ChannelI // Interface to native vehicle
	messages  chan messages.Message
	callbacks callback.Map
	parser    messages.Parser
	done      chan bool
}

// newSDKAdmin creates a new sdkVehicle
func newSDKChannel(i ChannelI, parser messages.Parser) *sdkChannel {
	ch := &sdkChannel{
		i:         i,
		messages:  make(chan messages.Message),
		callbacks: callback.NewMap(),
		parser:    parser,
		done:      make(chan bool),
	}

	i.SetListener(ch)

	return ch
}

func (ch *sdkChannel) Send(msg messages.Message) error {
	json, err := msg.ToJSON()

	if err != nil {
		return err
	}

	if msg.IsRequest() {
		ch.callbacks.Add(msg.ID, msg.Callback())
	}

	go ch.i.HandleMessage(json)
	return nil
}

func (ch *sdkChannel) Recv() <-chan messages.Message {
	return ch.messages
}

func (ch *sdkChannel) Done() <-chan bool {
	return ch.done
}

func (ch *sdkChannel) Disconnect() {
	// TODO
}

func (ch *sdkChannel) OnMessage(json string) {
	msg, err := ch.parser.Parse([]byte(json))
	if err != nil {
		log.Println("error parsing json message")
		return
	}

	if msg.Type == "reply" {
		ch.onReply(msg)
		return
	}

	if msg.IsRequest() {
		msg.Callback().Listen(func(result interface{}, err error) {
			data := libs.JSONObject{
				"id": msg.ID,
			}
			if err != nil {
				data["error"] = err.Error()
			} else {
				data["result"] = result
			}

			ch.Send(messages.New("reply", data))
		})
	}

	select {
	case ch.messages <- msg:
	case <-time.After(time.Second):
		log.Println("warning: unhandled message")
	}
}

func (ch *sdkChannel) onReply(msg messages.Message) {
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

	cb := ch.callbacks.Get(id)
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

func (ch *sdkChannel) Remove() {
	close(ch.done)
}
