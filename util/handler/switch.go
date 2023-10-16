package handler

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

type Switcher interface {
	Handle(i any) (any, error)
	RegisterInterface(handler any) error
	RegisterHandler(handler *Handler)
}

const (
	ErrRegisterInterface = lerr.Str("handler argument to RegisterInterface requires a func or a channel")
	ErrNoHandler         = lerr.Str("no handler found")
)

// Switch holds a set of Handlers and can invoke the correct handler by either
// name or type.
type Switch struct {
	handlers map[reflect.Type]*Handler
}

func NewSwitch(size int) *Switch {
	return &Switch{
		handlers: make(map[reflect.Type]*Handler, size),
	}
}

func Handlers(handlers ...any) (*Switch, error) {
	return NewSwitch(len(handlers)).RegisterInterfaces(handlers...)
}

// Handle will invoke a handler using 'i'. If 'i' is a string, it will try to
// match by name. Otherwise, it will try to match by type. This will result
// in unpredictable behavior if matching by name and type if one of the Handlers
// is a string handler.
func (s *Switch) Handle(i any) (any, error) {
	h, found := s.handlers[reflect.TypeOf(i)]

	if !found {
		return nil, ErrNoHandler
	}
	return h.Handle(i)
}

func (s *Switch) RegisterHandler(handler *Handler) {
	s.handlers[handler.Type()] = handler
}

// RegisterInterface a handler with ListenerMux. It will return the argument
// type for the handler. The handler must be either a handler function or a
// receiver channel.
func (s *Switch) RegisterInterface(handler any) error {
	v := reflect.ValueOf(handler)

	switch v.Kind() {
	case reflect.Func:
		return s.registerFunc(v)
	case reflect.Chan:
		return s.registerChan(v)
	}
	return ErrRegisterInterface
}

func (s *Switch) RegisterInterfaces(handlers ...any) (*Switch, error) {
	for _, h := range handlers {
		err := s.RegisterInterface(h)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}

func (s *Switch) registerFunc(v reflect.Value) error {
	h, err := ByValue(v)
	if err != nil {
		return err
	}

	s.RegisterHandler(h)
	return nil
}

func (s *Switch) registerChan(v reflect.Value) error {
	argType := v.Type().Elem()
	fn := func(i any) {
		v.Send(reflect.ValueOf(i))
	}
	s.handlers[argType] = &Handler{
		fn: reflect.ValueOf(fn),
	}

	return nil
}
