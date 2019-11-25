package bus

import (
	"errors"
	"reflect"
)

// ListenerMux takes objects off a channel and passes them into the handlers
// for that type.
type ListenerMux struct {
	in         <-chan interface{}
	handlers   map[reflect.Type][]reflect.Value
	ErrHandler func(error)
}

// NewListenerMux creates a ListenerMux for the in bus channel.
func NewListenerMux(in <-chan interface{}, errHandler func(error), handlers ...interface{}) (*ListenerMux, error) {
	lm := &ListenerMux{
		handlers:   make(map[reflect.Type][]reflect.Value),
		ErrHandler: errHandler,
		in:         in,
	}
	for _, h := range handlers {
		if _, err := lm.RegisterMuxHandler(h); err != nil {
			return nil, err
		}
	}
	return lm, nil
}

// SetIn sets the interface channel the ListerMux is listening on.
func (lm *ListenerMux) SetIn(in <-chan interface{}) {
	lm.in = in
}

// Run the ListenerMux
func (lm *ListenerMux) Run() {
	for i := range lm.in {
		lm.handle(i)
	}
}

func (lm *ListenerMux) handle(i interface{}) {
	v := []reflect.Value{reflect.ValueOf(i)}
	for _, l := range lm.handlers[v[0].Type()] {
		out := l.Call(v)
		if l := len(out); l > 0 {
			err, ok := out[l-1].Interface().(error)
			if ok && err != nil && lm.ErrHandler != nil {
				lm.ErrHandler(err)
			}
		}
	}
}

// RegisterMuxHandler a handler with ListenerMux. It will return the argument
// type for the handler.
func (lm *ListenerMux) RegisterMuxHandler(handler interface{}) (reflect.Type, error) {
	v := reflect.ValueOf(handler)

	switch v.Kind() {
	case reflect.Func:
		return lm.registerFunc(v)
	case reflect.Chan:
		return lm.registerChan(v)
	}
	return nil, errors.New("Register requires a func or a channel")
}

func (lm *ListenerMux) registerFunc(v reflect.Value) (reflect.Type, error) {
	t := v.Type()
	if t.NumIn() != 1 {
		return nil, errors.New("Can only register a func with exactly one argument: " + t.String())
	}
	argType := t.In(0)
	lm.handlers[argType] = append(lm.handlers[argType], v)
	return argType, nil
}

func (lm *ListenerMux) registerChan(v reflect.Value) (reflect.Type, error) {
	t := v.Type()
	argType := t.Elem()
	fn := func(i interface{}) {
		v.Send(reflect.ValueOf(i))
	}
	lm.handlers[argType] = append(lm.handlers[argType], reflect.ValueOf(fn))
	return argType, nil
}

// HandleError takes in an error and passes it to the ErrHandler if it is set.
func (lm *ListenerMux) HandleError(err error) {
	if lm.ErrHandler != nil {
		lm.ErrHandler(err)
	}
}

// SetErrorHandler to errHandler if it is not already defined.
func (lm *ListenerMux) SetErrorHandler(errHandler func(error)) {
	if lm.ErrHandler == nil {
		lm.ErrHandler = errHandler
	}
}
