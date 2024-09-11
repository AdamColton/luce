package timeout

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

var (
	wgType = reflector.Type[*sync.WaitGroup]()
)

const (
	// ErrTimeout is returned by After and Must when the time limit is reached.
	ErrTimeout = lerr.Str("timeout")

	// InvalidWaitMsg is the format of the error After returns when wait is not a
	// supported type. It is passed the type.
	InvalidWaitMsg = "expected wait to be function, got %s"
)

// Must calls After and will panic if an error is returned. The value of the
// panic is the error, for example ErrTimeout.
func Must(ms int, wait interface{}) {
	lerr.Panic(After(ms, wait))
}

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
// wait can be a *sync.WaitGroup, which must be a pointer because a copy of a
// WaitGroup does not behave correctly. After returns once its Wait returns.
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
	if v.Type() == wgType {
		return wg(d, wait.(*sync.WaitGroup))
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

func wg(d time.Duration, wg *sync.WaitGroup) (err error) {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()
	select {
	case <-time.After(d):
		err = ErrTimeout
	case <-ch:
	}
	return
}

// Signal is the element type of the channel that Run returns. Nothing is sent
// on it, the channel is closed.
type Signal struct{}

// Run calls fn in a Go routine and closes the returned channel when fn returns,
// so it can be passed to After to wait for fn with a time limit.
func Run(fn func()) <-chan Signal {
	ch := make(chan Signal)
	go func() {
		fn()
		close(ch)
	}()
	return ch
}
