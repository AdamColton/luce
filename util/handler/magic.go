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

// MethodsRegistrar finds the handler methods of a value by reflection. A method
// is a handler if it takes one argument and passes Handlers. Detailer describes
// the command that a method makes, and Log, if it is not nil, receives a line
// for each method that is registered. Handlers and Detailer must be set, as they
// are in DefaultRegistrar.
type MethodsRegistrar struct {
	Handlers filter.Filter[*reflector.Method]
	Detailer func(*reflector.Method) *CommandDetails
	Log      io.Writer
}

// HandlersFilter is the filter for the handler methods: those that take one
// argument and pass mr.Handlers.
func (mr MethodsRegistrar) HandlersFilter() filter.Filter[*reflector.Method] {
	return oneArg.And(mr.Handlers)
}

// DefaultRegistrar takes the methods whose names end in "Handler". The name of
// a command is the method's name without "Handler" and with the first letter in
// lower case. Its details come from the method with the same name but ending in
// "Usage", if there is one that takes no arguments and returns a
// *CommandDetails: StringHandler is described by StringUsage. A method named
// only "Handler" has no name to give a command, so it is disabled.
var DefaultRegistrar = MethodsRegistrar{
	Handlers: filter.MethodName(filter.Suffix("Handler")),
	Detailer: func(m *reflector.Method) *CommandDetails {
		ln := len(m.Name)
		if ln <= len("Handler") {
			return &CommandDetails{Disabled: true}
		}
		rs := []rune(m.Name[:ln-7]) // remove "Handler" suffix

		um := m.On.MethodByName(string(rs) + "Usage")
		var out *CommandDetails
		if um.Kind() != reflect.Invalid && usageMethodType(um.Type()) {
			out = um.Call(nil)[0].Interface().(*CommandDetails)
		} else {
			out = &CommandDetails{}
		}
		if out == nil {
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

// Register registers a Handler with s for each handler method of handlerType,
// which is a value so that the methods are bound to it. It returns the argument
// types of the Handlers that were registered and the errors from the methods
// that are not a valid Handler.
func (mr MethodsRegistrar) Register(s Switcher, handlerType any) ([]reflect.Type, error) {
	ms := reflector.MethodsOn(handlerType)
	handlers, _ := mr.HandlersFilter().SliceInPlace(ms)
	if mr.Log != nil {
		fmt.Fprintf(mr.Log, "On Type %s\n", reflect.TypeOf(handlerType))
	}
	return RegisterSwitchHandlerMethods(s, handlers.Iter(), mr.Log)
}

// RegisterSwitchHandlerMethods registers a Handler with s for each method the
// iterator returns, as Register does, and writes a line for each to log if it
// is not nil. A method that is not a valid Handler is skipped and its error is
// returned with the others.
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

// Commands returns a Command by name for each handler method of handlerType,
// leaving out the ones whose details are Disabled. The Action of a Command is
// its method, bound to handlerType. If two methods end up with the same name
// the last one is kept.
func (mr MethodsRegistrar) Commands(handlerType any) lmap.Wrapper[string, *Command] {
	out := lmap.New[string, *Command](nil)
	ms := reflector.MethodsOn(handlerType)
	handlersIter := mr.HandlersFilter().Iter(slice.NewIter(ms))
	handlersIter.For(func(m *reflector.Method) {
		cd := mr.Detailer(m)
		if !cd.Disabled {
			out.Set(cd.Name, &Command{
				Name:   cd.Name,
				Usage:  cd.Usage,
				Alias:  cd.Alias,
				Action: m.Func.Interface(),
			})
		}
	})
	return out
}

// CommandDetails describes the Command that a handler method makes. A method
// named <Name>Usage returns them for <Name>Handler when DefaultRegistrar is
// used.
type CommandDetails struct {
	// Name of the command. If it is empty, it comes from the method's name.
	Name string
	// Usage is a short description of the command.
	Usage string
	// Disabled leaves the command out.
	Disabled bool
	// Alias is a shorter name for the command.
	Alias string
}
