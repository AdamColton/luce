package channel

import (
	"time"
)

// Merge receives slices of data and merges them adding controlled delays to
// wait for more data.
type Merge[T any] struct {
	p             Pipe[[]T]
	MaxDelayMS    uint
	SingleDelayUS uint
	c             *Close
}

var (
	// MaxDelayMS is the default used when calling NewMerge
	MaxDelayMS uint = 1
	// SingleDelayUS is the default used when calling NewMerge
	SingleDelayUS uint = 10
)

// NewMerge creates an instance of Merge. The rcv and snd arguments are used to
// invoke NewPipe. MaxDelayMS and SingleDelayUS are set from the package level
// defaults.
func NewMerge[T any](rcv <-chan []T, snd chan<- []T) (m *Merge[T], retSnd chan<- []T, retRcv <-chan []T) {
	m = &Merge[T]{
		MaxDelayMS:    MaxDelayMS,
		SingleDelayUS: SingleDelayUS,
		c:             NewClose(),
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
	done := false
	max := time.NewTimer(time.Millisecond * time.Duration(m.MaxDelayMS))
	d := time.Microsecond * time.Duration(m.SingleDelayUS)
	single := time.NewTimer(d)
	for !done {
		select {
		case data := <-m.p.Rcv:
			buf = append(buf, data...)
			if !single.Stop() {
				<-single.C
			}
			single.Reset(d)
		case <-max.C:
			done = true
		case <-single.C:
			done = true
		case <-m.c.OnClose:
			done = true
		}
	}
	if !m.c.Closed() {
		m.p.Snd <- buf
	}
}
