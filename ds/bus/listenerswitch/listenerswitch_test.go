package listenerswitch_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/bus/listenerswitch"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestListenerMux(t *testing.T) {
	strCh := make(chan string)
	ch := make(chan any)
	done := make(chan bool)

	handler := func(s string) {
		strCh <- s
	}
	ls, err := listenerswitch.New(10, ch, nil, handler)
	assert.NoError(t, err)

	runner := func() {
		ls.Run()
		done <- true
	}
	testFn := func() {
		str := "test"
		ch <- str
		assert.Equal(t, "test", <-strCh)
		close(ch)
		assert.True(t, <-done)
	}

	go runner()
	assert.NoError(t, timeout.After(5, testFn))

	ch = make(chan any)
	ls.SetIn(ch)
	go runner()
	assert.NoError(t, timeout.After(5, testFn))

}

func TestHandle(t *testing.T) {
	errCh := make(chan error, 1)
	ls, err := listenerswitch.New(10, nil, errCh, func(i int) string { return strconv.Itoa(i * 2) })
	assert.NoError(t, err)

	out, err := ls.Handle(21)
	assert.NoError(t, err)
	assert.Equal(t, "42", out)
	assert.Len(t, errCh, 0)

	// The error goes to the error handler and back to the caller.
	out, err = ls.Handle("no handler for a string")
	assert.Nil(t, out)
	assert.Equal(t, handler.ErrNoHandler, err)
	assert.Equal(t, handler.ErrNoHandler, <-errCh)
}

func TestSetErrorHandler(t *testing.T) {
	ls, err := listenerswitch.New(10, nil, nil)
	assert.NoError(t, err)

	// With no error handler the error is only returned.
	_, err = ls.Handle(1)
	assert.Equal(t, handler.ErrNoHandler, err)

	var got []error
	assert.NoError(t, ls.SetErrorHandler(func(err error) { got = append(got, err) }))
	ls.Handle(1)
	assert.Equal(t, []error{handler.ErrNoHandler}, got)

	// A handler that can't be used leaves the current one in place.
	assert.Equal(t, lerr.ErrHandlerFunc, ls.SetErrorHandler("abc"))
	ls.Handle(2)
	assert.Len(t, got, 2)

	assert.NoError(t, ls.SetErrorHandler(nil))
	ls.Handle(3)
	assert.Len(t, got, 2)
}

func TestRunErrors(t *testing.T) {
	ch := make(chan any)
	errCh := make(chan error, 2)
	ls, err := listenerswitch.New(10, ch, errCh, func(s string) {})
	assert.NoError(t, err)

	go func() {
		ch <- 1
		ch <- "handled"
		ch <- 2.5
		close(ch)
	}()
	assert.NoError(t, timeout.After(1000, ls.Run))

	// Values without a handler are reported to the error handler.
	close(errCh)
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	assert.Equal(t, []error{handler.ErrNoHandler, handler.ErrNoHandler}, errs)
}

func TestListenerMuxErr(t *testing.T) {
	ls, err := listenerswitch.New(10, nil, nil, 123)
	assert.Nil(t, ls)
	assert.Equal(t, handler.ErrRegisterInterface, err)

	ls, err = listenerswitch.New(10, nil, "abc", 123)
	assert.Nil(t, ls)
	assert.Equal(t, lerr.ErrHandlerFunc, err)
}
