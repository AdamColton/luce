package channel

import (
	"time"

	"github.com/adamcolton/luce/lerr"
)

// ErrTimeout is returned by Timeout and TimeoutMS when nothing is received in
// time.
const ErrTimeout = lerr.Str("timeout")

// Timeout receives on ch for the duration d. If nothing is received in that
// time, the zero value and ErrTimeout are returned. A closed channel returns its
// zero value at once, with no error.
func Timeout[T any](d time.Duration, ch <-chan T) (t T, err error) {
	select {
	case t = <-ch:
	case <-time.After(d):
		err = ErrTimeout
	}
	return
}

// TimeoutMS is Timeout with the duration given in milliseconds.
func TimeoutMS[T any](ms int, ch <-chan T) (t T, err error) {
	d := time.Millisecond * time.Duration(ms)
	return Timeout(d, ch)
}
