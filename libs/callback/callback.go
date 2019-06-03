package callback

import (
	"errors"
	"sync"
	"time"
)

type Callback struct {
	sync.RWMutex
	done      chan bool
	res       interface{}
	err       error
	timeout   *time.Timer
	listeners []func(interface{}, error)
}

func New() *Callback {
	return &Callback{
		done: make(chan bool),
	}
}

func (cb *Callback) Timeout(timeout time.Duration) *Callback {
	cb.Lock()
	if cb.timeout == nil {
		cb.timeout = time.AfterFunc(timeout, func() {
			cb.Reject(errors.New("timeout"))
		})
	}
	cb.Unlock()

	return cb
}

func (cb *Callback) Wait() (interface{}, error) {
	<-cb.done
	return cb.res, cb.err
}

func (cb *Callback) Cancel() {
	cb.Reject(errors.New("canceled"))
}

func (cb *Callback) Resolve(data interface{}) bool {
	return cb.finish(data, nil)
}

func (cb *Callback) Reject(err error) bool {
	return cb.finish(nil, err)
}

func (cb *Callback) finish(res interface{}, err error) bool {
	cb.Lock()

	select {
	case <-cb.done:
		cb.Unlock()
		return false
	default:
	}

	if cb.timeout != nil {
		cb.timeout.Stop()
	}
	cb.res = res
	cb.err = err
	close(cb.done)

	listeners := cb.listeners
	cb.listeners = nil

	cb.Unlock()

	for _, fn := range listeners {
		fn(cb.res, cb.err)
	}

	return true
}

func (cb *Callback) Listen(fn func(interface{}, error)) {
	cb.Lock()
	select {
	case <-cb.done:
		fn(cb.res, cb.err)
	default:
		cb.listeners = append(cb.listeners, fn)
	}
	cb.Unlock()
}
