package messages

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
)

const (
	// REQUEST is a "verb" set to a message that awaits a response
	REQUEST = "req"
	// UPDATE is a "verb" set to a message that doesn't awaits a response
	UPDATE = "upd"
)

// Message represents a websocket message
type Message struct {
	ID   string             `json:"id"`
	Verb string             `json:"verb"`
	Type string             `json:"type"`
	Data interface{}        `json:"data"`
	cb   *callback.Callback `json:"-"`
}

// MessageType is a subset of Message used when only interested in it's type
type MessageType struct {
	Type string `json:"type"`
}

var inc uint64

// MessageID generates a unique id for a message
func id() string {
	//inc++
	atomic.AddUint64(&inc, 1)
	return strconv.FormatUint(inc, 10)
}

// New creates a new message
func New(t string, d interface{}) Message {
	return create(UPDATE, t, d)
}

// NewRequest creates a new request message (expects a response)
func NewRequest(t string, d interface{}, cb *callback.Callback) Message {
	msg := create(REQUEST, t, d)

	if cb == nil {
		cb = callback.New()
	}
	msg.cb = cb

	return msg
}

func create(verb string, t string, d interface{}) Message {
	return Message{
		ID:   id(),
		Verb: verb,
		Type: t,
		Data: d,
	}
}

// IsRequest checks if this message is a request
func (msg Message) IsRequest() bool {
	return msg.cb != nil
}

// ToRequest set the type of this message to "REQUEST" and sets a callback
func (msg Message) ToRequest(cb *callback.Callback) Message {
	msg.Verb = REQUEST
	msg.cb = cb
	return msg
}

//func (msg Message) callback(cb *callback.Callback) {
//	msg.cb = cb
//	if msg.cb == nil {
//		msg.verb = REQUEST
//	} else {
//		msg.verb = UPDATE
//	}
//}

// Callback returns the messages callback object
func (msg Message) Callback() *callback.Callback {
	return msg.cb
}

// SubType splits the type of the message by ":" and returns the element
// at the specified index
func (msg Message) SubType(i int) string {
	t := strings.Split(msg.Type, ":")
	if i < len(t) {
		return t[i]
	}

	return ""
}

// JSONData returns the message's data object as a libs.JSONObject
func (msg Message) JSONData() libs.JSONObject {
	if d, ok := msg.Data.(*libs.JSONObject); ok {
		return *d
	}

	return nil
}

// ToJSON returns the message as a json string
func (msg Message) ToJSON() (string, error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Reply sends a reply to a request message
func (msg Message) Reply(result interface{}, err error) error {
	if !msg.IsRequest() {
		return errors.New("not a request")
	}

	if err != nil {
		msg.cb.Reject(err)
	} else {
		msg.cb.Resolve(result)
	}

	return nil
}
