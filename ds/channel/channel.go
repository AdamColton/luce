package channel

import (
	"time"

	"github.com/adamcolton/luce/lerr"
)

// ErrTimeout is returned by Timeout
const ErrTimeout = lerr.Str("timeout")

// Timeout receives on ch for ms milliseconds. If nothing is received in that
// time, ErrTimeout is returned.
func Timeout[T any](ms int, ch <-chan T) (t T, err error) {
	d := time.Millisecond * time.Duration(ms)
	select {
	case t = <-ch:
	case <-time.After(d):
		err = ErrTimeout
	}
	return
}

// Slice places the contents of a slice on a channel. This is done as a blocking
// operation, however, if ch is nil a channel is created with a buffer large
// enough to hold the slice and it will not block, returning the populated
// channel.
func Slice[T any](s []T, ch chan<- T) (out chan T) {
	if ch == nil {
		out = make(chan T, len(s))
		ch = out
	}
	for _, t := range s {
		ch <- t
	}
	return
}
