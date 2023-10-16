package handler

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

// Switcher is the registering and invoking side of a Switch, so that something
// else can stand in for it. Handle invokes the Handler for the type of 'i'.
type Switcher interface {
	Handle(i any) (any, error)
	RegisterInterface(handler any) error
	RegisterHandler(handler *Handler)
}

const (
	// ErrRegisterInterface is returned by RegisterInterface for an argument that is
	// not a func or a channel.
	ErrRegisterInterface = lerr.Str("handler argument to RegisterInterface requires a func or a channel")
	// ErrNoHandler is returned by Handle when no Handler is registered for the
	// type it was given.
	ErrNoHandler = lerr.Str("no handler found")
)

// Switch holds a set of Handlers and can invoke the correct handler by the type
// of the argument. It is not safe for concurrent use.
type Switch struct {
	handlers map[reflect.Type]*Handler
}

// NewSwitch creates a Switch with room for size Handlers.
func NewSwitch(size int) *Switch {
	return &Switch{
		handlers: make(map[reflect.Type]*Handler, size),
	}
}

// Handlers creates a Switch and registers each of the handlers with
// RegisterInterfaces. It returns the Switch and the first error.
func Handlers(handlers ...any) (*Switch, error) {
	return NewSwitch(len(handlers)).RegisterInterfaces(handlers...)
}

// Handle will invoke the Handler registered for the type of 'i' and return its
// results. It returns ErrNoHandler if there is none. A nil 'i' invokes the
// Handler that takes no argument.
func (s *Switch) Handle(i any) (any, error) {
	h, found := s.handlers[reflect.TypeOf(i)]

	if !found {
		return nil, ErrNoHandler
	}
	return h.Handle(i)
}

// RegisterHandler adds the Handler, replacing any Handler for the same argument
// type. Handlers that take no argument all use the nil type, so only one can be
// registered, and it is the one Handle(nil) invokes.
func (s *Switch) RegisterHandler(handler *Handler) {
	s.handlers[handler.Type()] = handler
}

// RegisterInterface registers a handler, which must be either a func that
// Handler allows or a channel. A channel of T is registered as the Handler for
// T that sends the value on the channel, so Handle waits until it is received.
// It returns the error from ByValue for a bad func, and ErrRegisterInterface for
// anything else.
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

// RegisterInterfaces registers each of the handlers with RegisterInterface and
// stops at the first error. Those before it stay registered. It returns the
// Switch so calls can be chained.
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
