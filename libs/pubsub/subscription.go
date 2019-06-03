package pubsub

type Subscriber struct {
	inbox chan interface{}
	done  chan bool
}

func NewSubscriber() Subscriber {
	return Subscriber{
		inbox: make(chan interface{}),
		done:  make(chan bool),
	}
}

func (s Subscriber) Recv() <-chan interface{} {
	return s.inbox
}

func (s Subscriber) publish(data interface{}) {
	select {
	case s.inbox <- data:
	case <-s.done:
	}
}

func (s Subscriber) finish() {
	close(s.done)
}
