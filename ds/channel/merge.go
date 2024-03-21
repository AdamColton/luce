package channel

import (
	"time"

	"github.com/adamcolton/luce/util/dbltimer"
)

// Merge receives slices of data and merges them adding controlled delays to
// wait for more data.
type Merge[T any] struct {
	p           Pipe[[]T]
	MaxDelay    time.Duration
	SingleDelay time.Duration
	c           *Close
}

var (
	// MaxDelay is the default used when calling NewMerge
	MaxDelay = time.Millisecond
	// SingleDelay is the default used when calling NewMerge
	SingleDelay = 10 * time.Microsecond
)

// NewMerge creates an instance of Merge. The rcv and snd arguments are used to
// invoke NewPipe. MaxDelayMS and SingleDelayUS are set from the package level
// defaults.
func NewMerge[T any](rcv <-chan []T, snd chan<- []T) (m *Merge[T], retSnd chan<- []T, retRcv <-chan []T) {
	m = &Merge[T]{
		MaxDelay:    MaxDelay,
		SingleDelay: SingleDelay,
		c:           NewClose(),
	}
	m.p, retSnd, retRcv = NewPipe(rcv, snd)
	return
}

// Run invokes Cycle everytime data is received on In. This adds at least
// SingleDelayUS latentcy.
func (m *Merge[T]) Run() {
	for data := range m.p.Rcv {
		m.Cycle(data)
	}
	if m.c.Close() {
		close(m.p.Snd)
	}
}

// Cycle receives on In and combines all the slices it recevies and sends them
// to Out. It will receive for a maximum of MaxDelayMS or if it goes
// SingleDelayUS without receiving anything.
func (m *Merge[T]) Cycle(buf []T) {
	timerDone := make(chan bool)
	dt := dbltimer.New(m.MaxDelay, m.SingleDelay, func() {
		timerDone <- true
	})

	done := false
	for !done {
		select {
		case done = <-timerDone:
		case <-m.c.OnClose:
			done = true
		case data := <-m.p.Rcv:
			if data != nil {
				buf = append(buf, data...)
				dt.Reset()
			} else {
				// m.p.Rcv might be closed, this prevent the loop from running
				// continuously until the timer runs out and guarentees that
				// when the timer does run out done will be updated
				done, _ = Timeout(m.SingleDelay/2, timerDone)
			}
		}
	}
	if !m.c.Closed() {
		m.p.Snd <- buf
	}
}
