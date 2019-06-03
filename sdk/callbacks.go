package sdk

import (
	"log"
	"sync"
	"sync/atomic"
)

var callbacks sync.Map

type cbKey interface{}

var inc int64

func id() int64 {
	atomic.AddInt64(&inc, 1)
	return inc
}

func cb(fn func(err error)) int64 {
	id := id()
	callbacks.Store(cbKey(id), fn)
	return id
}

func callCb(id int64, err error) {
	val, ok := callbacks.Load(cbKey(id))
	if ok {
		cb, ok := val.(func(err error))
		if ok {
			cb(err)
		} else {
			log.Println("could not call callback")
		}
		callbacks.Delete(cbKey(id))
	} else {
		log.Println("callback does not exist")
	}
}
