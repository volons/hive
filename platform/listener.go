package platform

import (
	"github.com/volons/hive/messages"
)

type listener struct {
	msgID string
	ch    chan *messages.Message
	done  chan bool
}

func newListener(msgID string) listener {
	return listener{
		msgID: msgID,
		ch:    make(chan *messages.Message),
		done:  make(chan bool),
	}
}
