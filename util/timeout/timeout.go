package timeout

import (
	"fmt"
	"reflect"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

const (
	// ErrTimeout is returned by After and Must when the time limit is reached.
	ErrTimeout = lerr.Str("timeout")

	// InvalidWaitMsg is the format of the error After returns when wait is not a
	// supported type. It is passed the type.
	InvalidWaitMsg = "expected wait to be function, got %s"
)

// After waits for wait and returns ErrTimeout if it has not completed after ms
// milliseconds. It does not stop what it is waiting for: on a timeout the
// function keeps running in its Go routine.
//
// wait can be a function. It is called in a Go routine with no arguments, so it
// must not take any. If its last return value is an error and it returns in
// time, that error is returned.
//
// wait can be a channel. For a send only channel, the zero value of the
// channel's type is sent. For any other channel, After receives from it. If
// that does not happen in time, ErrTimeout is returned. If the value received is
// an error that is not nil, that error is returned.
//
// If wait is not a supported type, an error made from InvalidWaitMsg is
// returned. wait must not be nil.
func After(ms int, wait interface{}) error {
	d := time.Millisecond * time.Duration(ms)
	v := reflect.ValueOf(wait)
	switch v.Kind() {
	case reflect.Chan:
		if v.Type().ChanDir() == reflect.SendDir {
			return chSend(d, v)
		}
		return chRecv(d, v)
	case reflect.Func:
		return fn(d, v)
	}
	return fmt.Errorf(InvalidWaitMsg, v.Type())
}

func fn(d time.Duration, v reflect.Value) (err error) {
	ch := make(chan []reflect.Value)
	go func() {
		ch <- v.Call(nil)
	}()
	select {
	case <-time.After(d):
		err = ErrTimeout
	case out := <-ch:
		err = reflector.ReturnsErrCheck(out)
	}
	return
}

func chSend(d time.Duration, v reflect.Value) error {
	i, _, _ := reflect.Select([]reflect.SelectCase{
		{
			Dir:  reflect.SelectSend,
			Chan: v,
			Send: reflect.Zero(v.Type().Elem()),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(d)),
		},
	})
	if i == 1 {
		return ErrTimeout
	}
	return nil
}

func chRecv(d time.Duration, v reflect.Value) error {
	i, r, _ := reflect.Select([]reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: v,
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(time.After(d)),
		},
	})
	if i == 1 {
		return ErrTimeout
	}
	err, _ := r.Interface().(error)
	return err
}
