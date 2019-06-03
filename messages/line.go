package messages

import (
	"errors"
	"sync"

	"github.com/volons/hive/libs"
)

type Line struct {
	lock sync.RWMutex

	closeOnDisconnect bool

	receive chan Message
	peer    *Line

	done libs.Done
}

func NewLine(closeOnDisconnect bool) *Line {
	return &Line{
		closeOnDisconnect: closeOnDisconnect,
		receive:           make(chan Message),
		done:              libs.NewDone(),
	}
}

func (l *Line) Connection() <-chan interface{} {
	return nil
}

func (l *Line) Connected() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.peer != nil
}

func (l *Line) Connect(peer *Line) {
	l.Disconnect()
	peer.Disconnect()

	l.setPeer(peer)
	peer.setPeer(l)
}

func (l *Line) Disconnect() {
	l.lock.Lock()
	peer := l.peer
	l.peer = nil
	l.lock.Unlock()

	if peer != nil {
		peer.Disconnect()

		if l.closeOnDisconnect {
			l.Close()
		}
	}
}

func (l *Line) setPeer(peer *Line) {
	l.lock.Lock()
	l.peer = peer
	l.lock.Unlock()
}

//func (l *Line) Request(msg Message) *callback.Callback {
//	return l.RequestWithCallback(msg, callback.New())
//}
//
//func (l *Line) RequestWithCallback(msg Message, cb *callback.Callback) *callback.Callback {
//	msg.callback(cb)
//
//	err := l.Send(msg)
//	if err != nil {
//		cb.Reject(err)
//	}
//
//	return cb
//}

func (l *Line) Peer() *Line {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.peer
}

func (l *Line) Send(msg Message) error {
	peer := l.Peer()

	if peer == nil {
		return errors.New("not connected")
	}

	return peer.Push(msg)
}

func (l *Line) Push(msg Message) error {
	select {
	case l.receive <- msg:
	case <-l.Done():
		return errors.New("disconnected")
	}

	return nil
}

func (l *Line) Done() <-chan bool {
	return l.done.WaitCh()
}

func (l *Line) Close() {
	l.done.Done()
}

func (l *Line) Recv() <-chan Message {
	return l.receive
}
