package websocket

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync/atomic"

	"github.com/volons/hive/libs"
)

// Message represents a websocket message
type Message struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	IData interface{} `json:"data"`
}

var inc uint64

// MessageID generates a unique id for a message
func id() string {
	//inc++
	atomic.AddUint64(&inc, 1)
	return strconv.FormatUint(inc, 10)
}

// NewMessage creates a new message to send
func NewMessage(t string, d interface{}) Message {
	return Message{
		ID:    id(),
		Type:  t,
		IData: d,
	}
}

// ParseMessage parses a json string and
// creates a message object
func ParseMessage(message []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(message, &msg)
	return msg, err
}

// Data returns the message's data object as a libs.JSONObject
func (msg Message) Data() libs.JSONObject {
	if d, ok := msg.IData.(map[string]interface{}); ok {
		return libs.JSONObject(d)
	}

	return nil
}

// ToJSON returns the message as a json string
func (msg Message) ToJSON() (string, error) {
	return libs.ToJSON(msg)
}

// Result unwraps extracts the result and error data from a callback massage
// Callback example:
// { "type": "callback", "data": { "id": "<req msg id>", "result": { ... }, "error": "?" } }
func (msg Message) Result() (libs.JSONObject, error) {
	data := msg.Data()

	if data != nil {
		if err, ok := data.GetString("error"); ok && err != "" {
			return nil, errors.New(err)
		}
		if result, ok := data.GetObj("result"); ok {
			return result, nil
		}
	}

	return nil, nil
}
