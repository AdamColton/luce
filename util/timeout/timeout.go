// Package timeout is primarily intended as a test utility. It provides the
// function timeout.After(ms, wait) which will wait ms number of milliseconds
// for the wait object.
//
// The wait value can be a function. It will be called with no arguments. If it
// does not return before the timeout duration a TimeoutError is returned. If
// the last return value of the function is of type error and the function does
// not timeout but does return an error, that error will be returned.
//
// The wait value can be a channel. If it is a send only channel, the zero value
// for the channel will be sent. If it blocks longer than the timeout duration,
// a TimeoutError is returned. If the channel can receive, it will try for the
// timeout duration. If it does not receive within the duration, a TimeoutError
// is returned. If the channel does receive but the value that comes through the
// channel is an interface that fulfills error and is not nil, that error value
// is returned.
//
// The wait value can be a *sync.WaitGroup. It must be a pointer to a WaitGroup,
// passing in a WaitGroup by value causes it's Wait() method to not behave
// correctly.
//
// If the wait value is not a valid type an InvalidWait error is returned.
package timeout

import (
	"fmt"
	"reflect"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

const (
	// ErrTimeout is returned when a timeout occures
	ErrTimeout = lerr.Str("timeout")

	// InvalidWaitMsg returned by timeout.InvalidWaitMsg.Error()
	InvalidWaitMsg = "expected wait to be function, got %s"
)

// After returns Timeout when a specified number of milliseconds (ms) have
// passed if wait has not completed. If wait is not a valid type InvalidWait is
// returned.
func After(ms int, wait interface{}) error {
	d := time.Millisecond * time.Duration(ms)
	v := reflect.ValueOf(wait)
	switch v.Kind() {
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
