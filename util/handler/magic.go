package handler

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/reflector"
)

type MethodsRegistrar struct {
	Handlers filter.Filter[*reflector.Method]
}

func (mr MethodsRegistrar) HandlersFilter() filter.Filter[*reflector.Method] {
	return oneArg.And(mr.Handlers)
}

var DefaultRegistrar = MethodsRegistrar{
	Handlers: filter.MethodName(filter.Prefix("Handle")),
}

var (
	oneArg = filter.NumIn(filter.EQ(1)).Method()
)

func (mr MethodsRegistrar) Register(s Switcher, handlerType any) ([]reflect.Type, error) {
	ms := reflector.MethodsOn(handlerType)
	handlersIter := mr.HandlersFilter().Iter(slice.NewIter(ms))
	return RegisterSwitchHandlerMethods(s, handlersIter)
}

func RegisterSwitchHandlerMethods(s Switcher, handlersIter iter.Iter[*reflector.Method]) ([]reflect.Type, error) {
	var ts []reflect.Type
	var errs lerr.Many
	iter.Wrap(handlersIter).For(func(m *reflector.Method) {
		h, err := ByValue(m.Func, "")
		errs = errs.Add(err)
		if h != nil {
			s.RegisterHandler(h)
			ts = append(ts, h.Type())
		}
	})
	return ts, errs.Cast()
}
