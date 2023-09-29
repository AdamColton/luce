package bus

import (
	"reflect"

	"github.com/adamcolton/luce/util/handler"
)

// ListenerMethodsRegistrar finds the handler methods of a value, as
// handler.MethodsRegistrar does, and registers them on a Listener.
type ListenerMethodsRegistrar struct {
	handler.MethodsRegistrar
}

// DefaultRegistrar uses handler.DefaultRegistrar to find the handlers.
var DefaultRegistrar = ListenerMethodsRegistrar{handler.DefaultRegistrar}

// Register registers the handler methods of handlerType, a value that the methods
// are bound to, with the Listener, and then registers the argument type of each
// with the Listener's Receiver so that it can receive them. A method that is not
// a valid handler doesn't stop the others from being registered, but its error
// is returned. If registering a type fails Register stops and returns that
// error. The argument types must not be interfaces.
func (lmr ListenerMethodsRegistrar) Register(l Listener, handlerType any) error {
	ts, err := lmr.MethodsRegistrar.Register(l, handlerType)
	for _, t := range ts {
		i := reflect.New(t).Elem().Interface()
		if typeErr := l.RegisterType(i); typeErr != nil {
			return typeErr
		}
	}
	return err
}
