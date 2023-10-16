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

type key struct {
	reflect.Type
	name string
}

type Switch struct {
	handlers map[key]*Handler
}

func NewSwitch(size int) *Switch {
	return &Switch{
		handlers: make(map[key]*Handler, size),
	}
}

func (s *Switch) Handle(i any) (any, error) {
	var (
		h     *Handler
		found bool
	)

	if name, ok := i.(string); ok {
		h, found = s.handlers[key{
			name: name,
		}]
		if found {
			i = nil
		}
	}

	if !found {
		h, found = s.handlers[key{
			Type: reflect.TypeOf(i),
		}]
	}

	if !found {
		return nil, ErrNoHandler
	}
	return h.Handle(i)
}

func (s *Switch) RegisterHandler(handler *Handler) {
	s.handlers[handler.key()] = handler
}

// RegisterInterface a handler with ListenerMux. It will return the argument
// type for the handler.
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
	h, err := ByValue(v, "")
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
	k := key{
		Type: argType,
	}
	s.handlers[k] = &Handler{
		fn: reflect.ValueOf(fn),
	}

	return nil
}
