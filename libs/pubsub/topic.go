package pubsub

import "sync"

type Topic struct {
	sync.RWMutex
	subscribers []Subscriber
}

func NewTopic() *Topic {
	return &Topic{}
}

func (t *Topic) Subscription() Subscriber {
	s := NewSubscriber()

	t.Lock()
	t.subscribers = append(t.subscribers, s)
	t.Unlock()

	return s
}

func (t *Topic) Subscribe(s Subscriber) {
	t.Lock()
	t.subscribers = append(t.subscribers, s)
	t.Unlock()
}

func (t *Topic) Unsubscribe(s Subscriber) {
	t.Lock()
	defer t.Unlock()

	s.finish()

	for i, sub := range t.subscribers {
		if sub == s {
			t.subscribers = append(t.subscribers[:i], t.subscribers[i+1:]...)
			return
		}
	}
}

func (t *Topic) Publish(data interface{}) {
	t.RLock()
	defer t.RUnlock()

	for _, sub := range t.subscribers {
		sub.publish(data)
	}
}
