package handler

import (
	"reflect"
	"unicode"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/reflector"
)

type MethodsRegistrar struct {
	Handlers filter.Filter[*reflector.Method]
	Namer    func(*reflector.Method) (name, usage string, use bool)
}

func (mr MethodsRegistrar) HandlersFilter() filter.Filter[*reflector.Method] {
	return oneArg.And(mr.Handlers)
}

var DefaultRegistrar = MethodsRegistrar{
	Handlers: filter.MethodName(filter.Suffix("Handler")),
	Namer: func(m *reflector.Method) (name, usage string, use bool) {
		use = true
		ln := len(m.Name)
		rs := []rune(m.Name[:ln-7])

		u := ""
		um := m.On.MethodByName(string(rs) + "Usage")
		if um.Kind() == reflect.Func && usageMethodType(um.Type()) {
			out := um.Call(nil)
			u = out[0].Interface().(string)
			if len(out) == 2 {
				use = out[1].Interface().(bool)
			}
		}

		rs[0] = unicode.ToLower(rs[0])
		return string(rs), u, use
	},
}

var (
	oneArg          = filter.NumIn(filter.EQ(1)).Method()
	usageMethodType = filter.NumInEq(0).And(
		filter.OutType(0, reflector.Type[string]()),
	).And(
		filter.NumOut(filter.EQ(1)).Or(
			filter.OutType(1, reflector.Type[bool]()),
		),
	).Filter
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
		h, err := ByValue(m.Func)
		errs = errs.Add(err)
		if h != nil {
			s.RegisterHandler(h)
			ts = append(ts, h.Type())
		}
	})
	return ts, errs.Cast()
}

func (mr MethodsRegistrar) Commands(handlerType any) slice.Slice[Command] {
	var out []Command
	ms := reflector.MethodsOn(handlerType)
	handlersIter := mr.HandlersFilter().Iter(slice.NewIter(ms))
	handlersIter.For(func(m *reflector.Method) {
		n, u, use := mr.Namer(m)
		if use {
			out = append(out, Command{
				Name:   n,
				Usage:  u,
				Action: m.Func.Interface(),
			})
		}
	})
	return out
}
