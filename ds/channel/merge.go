package channel

import (
	"time"

	"github.com/adamcolton/luce/util/dbltimer"
)

// Merge receives slices of data and merges them adding controlled delays to
// wait for more data. The slices are received on the Pipe's Rcv and the merged
// slice is sent on its Snd.
type Merge[T any] struct {
	p Pipe[[]T]
	// MaxDelay is the longest a Cycle waits for more data, counted from the start
	// of the Cycle.
	MaxDelay time.Duration
	// SingleDelay is how long a Cycle waits with nothing received before it ends.
	// It starts again each time data is received.
	SingleDelay time.Duration
	c           *Close

	// newTimer and timeout are dbltimer.New and Timeout except in tests, which
	// replace them to control time.
	newTimer func(hard, soft time.Duration, callback func()) timer
	timeout  func(d time.Duration, ch <-chan bool) (bool, error)
}

// timer is the part of a dbltimer.DoubleTimer that Merge uses.
type timer interface {
	Reset() bool
}

func newDoubleTimer(hard, soft time.Duration, callback func()) timer {
	return dbltimer.New(hard, soft, callback)
}

var (
	// MaxDelay is the default used when calling NewMerge
	MaxDelay = time.Millisecond
	// SingleDelay is the default used when calling NewMerge
	SingleDelay = 10 * time.Microsecond
)

// NewMerge creates an instance of Merge. The rcv and snd arguments are used to
// invoke NewPipe. MaxDelay and SingleDelay are set from the package level
// defaults. The Merge owns snd: Run closes it once rcv is closed.
func NewMerge[T any](rcv <-chan []T, snd chan<- []T) (m *Merge[T], retSnd chan<- []T, retRcv <-chan []T) {
	m = &Merge[T]{
		MaxDelay:    MaxDelay,
		SingleDelay: SingleDelay,
		c:           NewClose(),
		newTimer:    newDoubleTimer,
		timeout:     Timeout[bool],
	}
	m.p, retSnd, retRcv = NewPipe(rcv, snd)
	return
}

// Run invokes Cycle every time data is received on Rcv, until Rcv is closed.
// Then it closes Snd, so running it again returns at once. This adds at least
// SingleDelay of latency.
func (m *Merge[T]) Run() {
	for data := range m.p.Rcv {
		m.Cycle(data)
	}
	if m.c.Close() {
		close(m.p.Snd)
	}
}

// Cycle receives on Rcv, adds all the slices it receives to buf and sends the
// result on Snd. It will receive for a maximum of MaxDelay, or until it goes
// SingleDelay without receiving anything. Data is appended to buf, so buf may
// be modified. If Run has finished, Cycle returns without sending anything.
func (m *Merge[T]) Cycle(buf []T) {
	timerDone := make(chan bool)
	dt := m.newTimer(m.MaxDelay, m.SingleDelay, func() {
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
				// continuously until the timer runs out and guarantees that
				// when the timer does run out done will be updated
				done, _ = m.timeout(m.SingleDelay/2, timerDone)
			}
		}
	}
	if !m.c.Closed() {
		m.p.Snd <- buf
	}
}
