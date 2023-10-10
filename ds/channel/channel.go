package channel

import (
	"time"

	"github.com/adamcolton/luce/lerr"
)

// ErrTimeout is returned by Timeout
const ErrTimeout = lerr.Str("timeout")

// Timeout receives on ch for ms milliseconds. If nothing is received in that
// time, ErrTimeout is returned.
func Timeout[T any](d time.Duration, ch <-chan T) (t T, err error) {
	select {
	case t = <-ch:
	case <-time.After(d):
		err = ErrTimeout
	}
	return
}

func TimeoutMS[T any](ms int, ch <-chan T) (t T, err error) {
	d := time.Millisecond * time.Duration(ms)
	return Timeout(d, ch)
}
