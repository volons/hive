package callback

import (
	"sync"
	"time"
)

type Map struct {
	*sync.RWMutex
	data map[string]*Callback
}

func NewMap() Map {
	m := Map{
		RWMutex: &sync.RWMutex{},
		data:    make(map[string]*Callback),
	}

	return m
}

func (m Map) Add(id string, cb *Callback) bool {
	if cb == nil {
		cb = New()
	}

	m.Lock()

	if _, ok := m.data[id]; ok {
		m.Unlock()
		return false
	}

	m.data[id] = cb

	m.Unlock()

	cb.Timeout(time.Minute)
	cb.Listen(func(interface{}, error) {
		m.Lock()
		delete(m.data, id)
		m.Unlock()
	})

	return true
}

func (m Map) Get(id string) *Callback {
	m.RLock()
	defer m.RUnlock()

	return m.data[id]
}
