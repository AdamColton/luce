package bus

import (
	"reflect"
)

type listener struct {
	Receiver
	ListenerMuxer
	ch chan interface{}
}

// NewListener creates a Listener from a Receiver and a ListenerMuxer. It
// connects the Out channel on the Receiver to the In channel on the
// ListenerMuxer. If ListenerMuxer is nil, one will be created. If errHandler is
// not null it will be set on both the Receiver and ListenerMuxer.
func NewListener(r Receiver, lm ListenerMuxer, errHandler func(error), handlers ...interface{}) (Listener, error) {
	ch := make(chan interface{})
	r.SetOut(ch)
	if lm == nil {
		var err error
		lm, err = NewListenerMux(ch, errHandler)
		if err != nil {
			return nil, err
		}
	} else {
		lm.SetIn(ch)
		lm.SetErrorHandler(errHandler)
	}
	r.SetErrorHandler(errHandler)
	l := &listener{
		Receiver:      r,
		ListenerMuxer: lm,
		ch:            ch,
	}
	err := l.RegisterHandlers(handlers...)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Run both the Receiver and the ListenerMuxer.
func (l *listener) Run() {
	go func() {
		l.Receiver.Run()
		close(l.ch)
	}()
	l.ListenerMuxer.Run()
}

// RegisterHandler with the underlying ListenerMuxer and the argument to the
// handler is registered with the Receiver.
func (l *listener) RegisterHandlers(handlers ...interface{}) error {
	for _, h := range handlers {
		t, err := l.RegisterMuxHandler(h)
		if err != nil {
			return err
		}
		zeroVal := reflect.New(t).Elem().Interface()
		err = l.RegisterType(zeroVal)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetErrorHandler on both the underlying ListenerMuxer and Receiver.
func (l *listener) SetErrorHandler(errHandler func(error)) {
	l.ListenerMuxer.SetErrorHandler(errHandler)
	l.Receiver.SetErrorHandler(errHandler)
}
