package messages

import (
	"encoding/json"

	"github.com/volons/hive/libs/callback"
)

type typeGetterT func(string) interface{}

// Parser parses a message with the correct appropriate struct
// according to it's type
type Parser struct {
	typeGetter typeGetterT
}

// NewParser creates a new Parser instance with a function that returns
// a struct type for a given message type
func NewParser(typeGetter typeGetterT) Parser {
	return Parser{
		typeGetter: typeGetter,
	}
}

// Parse parses a json string and creates a message object
func (p Parser) Parse(data []byte) (Message, error) {
	var typ MessageType
	err := json.Unmarshal(data, &typ)
	if err != nil {
		return Message{}, err
	}

	var msg Message
	msg.Data = p.typeGetter(typ.Type)
	err = json.Unmarshal(data, &msg)
	if err != nil {
		return Message{}, err
	}

	if msg.Verb == REQUEST {
		msg.cb = callback.New()
	}

	return msg, nil
}
