package channel

import (
	"time"

	"github.com/adamcolton/luce/lerr"
)

const ErrTimeout = lerr.Str("timeout")

func Timeout[T any](ms int, ch <-chan T) (t T, err error) {
	d := time.Millisecond * time.Duration(ms)
	select {
	case t = <-ch:
	case <-time.After(d):
		err = ErrTimeout
	}
	return
}

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
