// Package listenerswitch routes values from a channel to the handler registered
// for their type.
package listenerswitch

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
)

// ListenerSwitch takes values off a channel and passes each to the handler
// registered for its type, as a handler.Switch does. Errors, including that
// there is no handler for a type, go to the error handler. It fulfills
// bus.ListenerSwitcher. SetIn and SetErrorHandler can't be used while it is
// running.
type ListenerSwitch struct {
	in <-chan any
	*handler.Switch
	lerr.ErrHandler
}

// New creates a ListenerSwitch with room for size handlers that takes values from
// the in channel. The errHandler is a func(error) or a channel of errors, as
// lerr.HandlerFunc accepts, or nil to ignore errors. Each of handlers is
// registered with RegisterInterface. It returns the first error.
func New(size int, in <-chan any, errHandler any, handlers ...any) (*ListenerSwitch, error) {
	ls := &ListenerSwitch{
		Switch: handler.NewSwitch(size),
		in:     in,
	}
	var err error
	ls.ErrHandler, err = lerr.HandlerFunc(errHandler)
	if err != nil {
		return nil, err
	}
	for _, h := range handlers {
		if err := ls.RegisterInterface(h); err != nil {
			return nil, err
		}
	}
	return ls, nil
}

// Handle passes i to the handler for its type, as Switch.Handle does. If that
// returns an error it is also passed to the error handler.
func (ls *ListenerSwitch) Handle(i any) (any, error) {
	out, err := ls.Switch.Handle(i)
	ls.ErrHandler.Handle(err)
	return out, err
}

// SetErrorHandler sets the error handler, a func(error) or a channel of errors,
// as lerr.HandlerFunc accepts, or nil to ignore errors. If it returns an error
// the error handler is unchanged.
func (ls *ListenerSwitch) SetErrorHandler(i any) error {
	h, err := lerr.HandlerFunc(i)
	if err != nil {
		return err
	}
	ls.ErrHandler = h
	return nil
}

// Run takes values off the channel and handles each of them. It returns when the
// channel is closed, and never returns if the channel is nil. A handler that
// panics stops Run.
func (ls *ListenerSwitch) Run() {
	for i := range ls.in {
		_, err := ls.Switch.Handle(i)
		ls.ErrHandler.Handle(err)
	}
}

// SetIn sets the channel that Run takes values from.
func (ls *ListenerSwitch) SetIn(in <-chan interface{}) {
	ls.in = in
}
