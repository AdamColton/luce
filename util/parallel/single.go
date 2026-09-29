package parallel

import "sync"

// Single guarantees that only one invocation of a function is running at a
// time. Calls to Run made while the function is already running return
// immediately without calling it; they are not queued.
type Single struct {
	// Mux guards the running state.
	Mux     sync.Mutex
	running bool
	fn      func()
}

// NewSingle creates a Single that will invoke fn.
func NewSingle(fn func()) *Single {
	return &Single{
		fn: fn,
	}
}

// Run invokes fn and blocks until it returns. If fn is already being run by
// another call to Run, it returns immediately without invoking fn. If fn
// panics, the Single will remain marked as running and fn will not be invoked
// again.
func (s *Single) Run() {
	s.Mux.Lock()
	if s.running {
		s.Mux.Unlock()
		return
	}
	s.running = true
	s.Mux.Unlock()

	s.fn()

	s.Mux.Lock()
	s.running = false
	s.Mux.Unlock()
}

// Running reports if fn is currently being run. It does not lock Mux, so it
// should only be relied on when there is some other synchronization with the
// call to Run.
func (s *Single) Running() bool {
	return s.running
}
