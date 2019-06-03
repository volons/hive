package messages

import "github.com/volons/hive/libs"

type Peer struct {
	libs.Done
	send    chan Message
	receive chan Message
}

func NewPeer() Peer {
	return Peer{
		send:    make(chan Message),
		receive: make(chan Message),
		Done:    libs.NewDone(),
	}
}
