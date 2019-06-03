package libs

import "sync"

// Done can be used to broadcast safely a done message
// across multiple goroutines
type Done struct {
	done chan bool
	once *sync.Once
}

// NewDone creates a new done struct
func NewDone() Done {
	return Done{
		done: make(chan bool),
		once: &sync.Once{},
	}
}

// Done closes the done channel (safe to call multiple times)
func (c Done) Done() {
	c.once.Do(func() {
		close(c.done)
	})
}

// Wait returns once the Done function is called
func (c Done) Wait() {
	<-c.done
}

// WaitCh returns a channel that will be closed
// when the done function is called
func (c Done) WaitCh() <-chan bool {
	return c.done
}
