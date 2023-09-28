package listener_test

import (
	"errors"
	"testing"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listener"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// mockReceiver is a Receiver that runs until it is stopped.
type mockReceiver struct {
	out        chan<- any
	stop       chan bool
	types      []any
	errHandler any
	typeErr    error
	handlerErr error
}

func newMockReceiver() *mockReceiver {
	return &mockReceiver{stop: make(chan bool)}
}

func (mr *mockReceiver) Run() {
	<-mr.stop
}

func (mr *mockReceiver) RegisterType(zeroValue any) error {
	mr.types = append(mr.types, zeroValue)
	return mr.typeErr
}

func (mr *mockReceiver) SetOut(out chan<- any) {
	mr.out = out
}

func (mr *mockReceiver) SetErrorHandler(errHandler any) error {
	mr.errHandler = errHandler
	return mr.handlerErr
}

// A Listener is both a Receiver and a ListenerSwitcher.
var (
	_ bus.Receiver         = (*listener.Listener)(nil)
	_ bus.ListenerSwitcher = (*listener.Listener)(nil)
)

func TestListener(t *testing.T) {
	r := newMockReceiver()
	errCh := make(chan error, 1)
	strCh := make(chan string)
	handler := func(str string) {
		strCh <- str
	}
	l, err := listener.New(10, r, errCh, handler)
	assert.NoError(t, err)

	// The type the handler takes is registered with the Receiver.
	assert.Equal(t, []any{""}, r.types)

	done := make(chan bool)
	go func() {
		l.Run()
		close(done)
	}()

	// What the Receiver sends goes to the handler for its type.
	assert.NoError(t, timeout.After(1000, func() {
		r.out <- "test"
		assert.Equal(t, "test", <-strCh)

		// A value with no handler goes to the error handler.
		r.out <- 5
		assert.Error(t, <-errCh)
	}))

	// Run returns when the Receiver has stopped.
	close(r.stop)
	assert.NoError(t, timeout.After(1000, done))
}

func TestNewSetsErrorHandlers(t *testing.T) {
	r := newMockReceiver()
	var errs []error
	errHandler := func(err error) { errs = append(errs, err) }

	l, err := listener.New(10, r, errHandler)
	assert.NoError(t, err)
	assert.NotNil(t, r.errHandler)

	// Without an error handler the Receiver is left alone.
	r = newMockReceiver()
	_, err = listener.New(10, r, nil)
	assert.NoError(t, err)
	assert.Nil(t, r.errHandler)

	// Both can be changed later.
	r = newMockReceiver()
	l, err = listener.New(10, r, nil)
	assert.NoError(t, err)
	assert.NoError(t, l.SetErrorHandler(errHandler))
	assert.NotNil(t, r.errHandler)
	_, err = l.Handle(5)
	assert.Error(t, err)
	assert.Len(t, errs, 1)
}

func TestSetErrorHandlerErrors(t *testing.T) {
	l, err := listener.New(10, newMockReceiver(), nil)
	assert.NoError(t, err)
	assert.Equal(t, lerr.ErrHandlerFunc, l.SetErrorHandler("not an error handler"))

	r := newMockReceiver()
	l, err = listener.New(10, r, nil)
	assert.NoError(t, err)
	r.handlerErr = errors.New("receiver refused")
	assert.EqualError(t, l.SetErrorHandler(func(error) {}), "receiver refused")
}

func TestNewErrors(t *testing.T) {
	// An error handler that can't be used.
	l, err := listener.New(10, newMockReceiver(), "not an error handler")
	assert.Nil(t, l)
	assert.Equal(t, lerr.ErrHandlerFunc, err)

	// The Receiver can refuse it.
	r := newMockReceiver()
	r.handlerErr = errors.New("receiver refused")
	l, err = listener.New(10, r, func(error) {})
	assert.Nil(t, l)
	assert.EqualError(t, err, "receiver refused")

	// A handler that can't be registered.
	l, err = listener.New(10, newMockReceiver(), nil, "not a handler")
	assert.Nil(t, l)
	assert.Error(t, err)
}

func TestRegisterHandlers(t *testing.T) {
	r := newMockReceiver()
	l, err := listener.New(10, r, nil)
	assert.NoError(t, err)

	assert.NoError(t, l.RegisterHandlers(func(s string) {}, func(i int) int { return i }))
	assert.Equal(t, []any{"", 0}, r.types)

	// A handler with no argument has no type to receive.
	err = l.RegisterHandlers(func() {})
	assert.EqualError(t, err, "a handler must take an argument, to give the type it handles")

	// It stops at the first error, and the ones before it stay registered.
	r.types = nil
	err = l.RegisterHandlers(func(f float64) {}, 5, func(b bool) {})
	assert.Error(t, err)
	assert.Equal(t, []any{0.0}, r.types)

	// The Receiver can refuse a type.
	r.typeErr = errors.New("unknown type")
	assert.EqualError(t, l.RegisterHandlers(func(b bool) {}), "unknown type")
}
