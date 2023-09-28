package listener

import (
	"reflect"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listenerswitch"
	"github.com/adamcolton/luce/util/handler"
)

type Listener struct {
	bus.Receiver
	bus.ListenerSwitcher
	ch chan interface{}
}

func New(muxSize int, r bus.Receiver, errHandler any, handlers ...interface{}) (*Listener, error) {
	ls, err := listenerswitch.New(muxSize, nil, errHandler)
	if err != nil {
		return nil, err
	}
	return FromMux(r, ls, handlers...)
}

// New creates a Listener from a Receiver and a ListenerSwitcher. It connects the
// Out channel on the Receiver to the In channel on the ListenerSwitcher. If
// ListenerSwitcher is nil, one will be created. If errHandler is not nil it will
// be set on both the Receiver and ListenerSwitcher.
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

// Run both the Receiver and the ListenerSwitcher.
func (l *Listener) Run() {
	go func() {
		l.Receiver.Run()
		close(l.ch)
	}()
	l.ListenerSwitcher.Run()
}

// RegisterHandler with the underlying ListenerSwitcher and the argument to the
// handler is registered with the Receiver.
func (l *Listener) RegisterHandlers(handlers ...any) error {
	for _, i := range handlers {
		h, err := handler.New(i)
		if err != nil {
			return err
		}
		l.RegisterHandler(h)
		t := h.Type()
		zeroVal := reflect.New(t).Elem().Interface()
		err = l.RegisterType(zeroVal)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetErrorHandler on both the underlying ListenerSwitcher and Receiver.
func (l *Listener) SetErrorHandler(errHandler any) error {
	err := l.ListenerSwitcher.SetErrorHandler(errHandler)
	if err != nil {
		return err
	}
	return l.Receiver.SetErrorHandler(errHandler)
}
