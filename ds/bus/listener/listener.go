// Package listener connects a bus receiver to a listener switch, so that the
// messages that arrive on a bus are passed to the handler for their type.
package listener

import (
	"reflect"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listenerswitch"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
)

// Listener takes data off a bus and passes each value to the handler for its
// type. The Receiver turns the data into values, which it sends to the
// ListenerSwitcher to be routed. It fulfills bus.Listener.
type Listener struct {
	bus.Receiver
	bus.ListenerSwitcher
	ch chan interface{}
}

// New creates a Listener for the Receiver, with a new ListenerSwitch that has
// room for muxSize handlers. If errHandler is not nil, it is set on both the
// Receiver and the ListenerSwitch. Each of the handlers is registered as
// RegisterHandlers does, and New returns the first error.
func New(muxSize int, r bus.Receiver, errHandler any, handlers ...interface{}) (*Listener, error) {
	ls, err := listenerswitch.New(muxSize, nil, errHandler)
	if err != nil {
		return nil, err
	}
	if errHandler != nil {
		err = r.SetErrorHandler(errHandler)
		if err != nil {
			return nil, err
		}
	}
	return FromMux(r, ls, handlers...)
}

// FromMux creates a Listener from a Receiver and a ListenerSwitcher, and neither
// can be nil. It connects the out channel of the Receiver to the in channel of
// the ListenerSwitcher, and registers the handlers as RegisterHandlers does.
func FromMux(r bus.Receiver, lm bus.ListenerSwitcher, handlers ...interface{}) (*Listener, error) {
	ch := make(chan any)
	lm.SetIn(ch)
	r.SetOut(ch)

	l := &Listener{
		Receiver:         r,
		ListenerSwitcher: lm,
		ch:               ch,
	}
	err := l.RegisterHandlers(handlers...)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Run runs the Receiver and the ListenerSwitcher, and returns when the Receiver
// has stopped and the ListenerSwitcher has handled everything it sent.
func (l *Listener) Run() {
	go func() {
		l.Receiver.Run()
		close(l.ch)
	}()
	l.ListenerSwitcher.Run()
}

// RegisterHandlers registers each handler, a func that handler.New accepts, with
// the ListenerSwitcher, and registers the type of its argument with the
// Receiver, so that the Receiver can receive it. A handler has to take an
// argument, and that must not be an interface type. It stops at the first error.
func (l *Listener) RegisterHandlers(handlers ...any) error {
	for _, i := range handlers {
		h, err := handler.New(i)
		if err != nil {
			return err
		}
		t := h.Type()
		if t == nil {
			return lerr.Str("a handler must take an argument, to give the type it handles")
		}
		l.RegisterHandler(h)
		zeroVal := reflect.New(t).Elem().Interface()
		err = l.RegisterType(zeroVal)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetErrorHandler sets the error handler on both the ListenerSwitcher and the
// Receiver, and stops at the first error.
func (l *Listener) SetErrorHandler(errHandler any) error {
	err := l.ListenerSwitcher.SetErrorHandler(errHandler)
	if err != nil {
		return err
	}
	return l.Receiver.SetErrorHandler(errHandler)
}
