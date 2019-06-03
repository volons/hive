package libs

import "sync/atomic"

// AtomicBool implements a thread safe boolean
type AtomicBool struct {
	flag int32
}

// Set sets the boolean value
func (b *AtomicBool) Set(val bool) {
	var i int32
	if val {
		i = 1
	}

	atomic.StoreInt32(&b.flag, int32(i))
}

// Get returns the boolean value
func (b *AtomicBool) Get() bool {
	if atomic.LoadInt32(&b.flag) != 0 {
		return true
	}

	return false
}

// Swap sets a new value and returns the old value
func (b *AtomicBool) Swap(val bool) bool {
	var i int32
	if val {
		i = 1
	}

	if atomic.SwapInt32(&b.flag, i) != 0 {
		return true
	}

	return false
}
