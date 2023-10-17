package handler

import (
	"fmt"
	"io"
	"reflect"
	"unicode"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/reflector"
)

type MethodsRegistrar struct {
	Handlers filter.Filter[*reflector.Method]
	Detailer func(*reflector.Method) *CommandDetails
	Log      io.Writer
}

func (mr MethodsRegistrar) HandlersFilter() filter.Filter[*reflector.Method] {
	return oneArg.And(mr.Handlers)
}

var DefaultRegistrar = MethodsRegistrar{
	Handlers: filter.MethodName(filter.Suffix("Handler")),
	Detailer: func(m *reflector.Method) *CommandDetails {
		ln := len(m.Name)
		rs := []rune(m.Name[:ln-7]) // remove "Handler" suffix

		um := m.On.MethodByName(string(rs) + "Usage")
		var out *CommandDetails
		if um.Kind() != reflect.Invalid && usageMethodType(um.Type()) {
			out = um.Call(nil)[0].Interface().(*CommandDetails)
		} else {
			out = &CommandDetails{}
		}
		rs[0] = unicode.ToLower(rs[0])
		if out.Name == "" {
			out.Name = string(rs)
		}

		return out
	},
}

var (
	oneArg          = filter.NumIn(filter.EQ(1)).Method()
	usageMethodType = filter.NumInEq(0).And(
		filter.OutType(0, reflector.Type[*CommandDetails]()),
	).Filter
)

func (mr MethodsRegistrar) Register(s Switcher, handlerType any) ([]reflect.Type, error) {
	ms := reflector.MethodsOn(handlerType)
	handlersIter := mr.HandlersFilter().Slice(ms).Iter()
	if mr.Log != nil {
		fmt.Fprintf(mr.Log, "On Type %s\n", reflect.TypeOf(handlerType))
	}
	return RegisterSwitchHandlerMethods(s, handlersIter, mr.Log)
}

func RegisterSwitchHandlerMethods(s Switcher, handlersIter liter.Iter[*reflector.Method], log io.Writer) ([]reflect.Type, error) {
	var ts []reflect.Type
	var errs lerr.Many
	liter.Wrap(handlersIter).For(func(m *reflector.Method) {
		h, err := ByValue(m.Func)
		errs = errs.Add(err)
		if h != nil {
			s.RegisterHandler(h)
			ts = append(ts, h.Type())
		}
		if log != nil {
			fmt.Fprintf(log, "%s %s\n", m.Name, m.Func.String())
		}
	})
	return ts, errs.Cast()
}

func (mr MethodsRegistrar) Commands(handlerType any) lmap.Map[string, *Command] {
	out := make(lmap.Map[string, *Command])
	ms := reflector.MethodsOn(handlerType)
	handlersIter := mr.HandlersFilter().Iter(slice.NewIter(ms))
	handlersIter.For(func(m *reflector.Method) {
		cd := mr.Detailer(m)
		if !cd.Disabled {
			out[cd.Name] = &Command{
				Name:   cd.Name,
				Usage:  cd.Usage,
				Alias:  cd.Alias,
				Action: m.Func.Interface(),
			}
		}
	})
	return out
}

type CommandDetails struct {
	Name     string
	Usage    string
	Disabled bool
	Alias    string
}
