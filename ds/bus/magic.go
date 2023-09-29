package bus

import (
	"reflect"

	"github.com/adamcolton/luce/util/handler"
)

type ListenerMethodsRegistrar struct {
	handler.MethodsRegistrar
}

var DefaultRegistrar = ListenerMethodsRegistrar{handler.DefaultRegistrar}

// RegisterHandlerType is a bit of reflection magic. Calls
// RegisterMuxHandlerType on the underlying ListenerMuxer and registers the
// returned types with the underlying Receiver.
func (lmr ListenerMethodsRegistrar) Register(l Listener, handlerType any) error {
	ts, err := lmr.MethodsRegistrar.Register(l, handlerType)
	if err != nil {
		return err
	}
	for _, t := range ts {
		i := reflect.New(t).Elem().Interface()
		err = l.RegisterType(i)
		if err != nil {
			return err
		}
	}
	return nil
}
