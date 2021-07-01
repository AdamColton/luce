package bus

import (
	"reflect"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
)

// RegisterMuxHandlerType is a bit of reflection magic. It takes a object and
// iterates over it's methods. Any methods that start with "Handler" will be
// registered with the ListenerMuxer. If there is a method named "ErrHandler"
// and the ListenerMuxer's ErrHandler field is nil, the field will be set to the
// method. A slice containing the arugments types of the handlers is returned.
func RegisterMuxHandlerType(lm ListenerMuxer, handlerType interface{}) ([]reflect.Type, error) {
	ms := reflector.ArgCount(filter.EQ.Int(1)).MethodFilter().
		And(reflector.MethodName(filter.Prefix("Handle"))).
		On(handlerType).Funcs()

	var ts []reflect.Type
	for _, m := range ms {
		at, err := lm.RegisterMuxHandler(m)
		if err != nil {
			return ts, err
		}
		ts = append(ts, at)
	}

	var errHandler func(err error)
	ok := reflector.
		MethodName(filter.EQ.String("ErrHandler")).
		One(handlerType).
		SetTo(&errHandler)
	if ok {
		lm.SetErrorHandler(errHandler)
	}

	return ts, nil
}

// RegisterHandlerType is a bit of reflection magic. Calls
// RegisterMuxHandlerType on the underlying ListenerMuxer and registers the
// returned types with the underlying Receiver.
func RegisterHandlerType(l Listener, handlerType interface{}) error {
	ts, err := RegisterMuxHandlerType(l, handlerType)
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
