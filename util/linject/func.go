package linject

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// FuncType is a helper when invoking a FuncInitilizer. All that is necessisary
// for a FuncInitilizer is the function type but so much logic is dedicated to
// inspecting the Target that it's convenient to provide that seperatly.
type FuncType interface {
	// Fn is the function type
	Fn() reflect.Type
	// DataType is a helper. The last argument of the function is a pointer
	// to a struct, DataType is that struct.
	Target() reflect.Type

	private()
}

// fnTypes privately backs the Func interface.
type fnTypes struct {
	target reflect.Type
	t      reflect.Type
}

func (ft *fnTypes) Target() reflect.Type {
	return ft.target
}

func (ft *fnTypes) Fn() reflect.Type {
	return ft.t
}

func (ft *fnTypes) private() {
}

// Initilizer represents a single injector. It inspects the FuncType to
// correctly apply the Injector. If the injection logic does not apply to the
// particular FuncType then nil is returned and the injector is not used with
// the given function. The injector may also use this initilization stage to
// dynamically change the logic in the returned FuncInjector rather than
// performing repetitive inspections on each invocation.
type Initilizer interface {
	Initilize(FuncType) Injector
}

// Injector holds injection logic that is invoked before a function is
// called. The args argument holds the values that will be passed into the
// function that is the subject of the injection. If a callback is returned,
// then the functions return values will be passed into the callback.
type Injector interface {
	Inject(args []reflect.Value) (callback func(rets []reflect.Value), err error)
}

// FuncInitilizers allows a slice of FuncInitilizer to be applied to a function.
// All FuncInitilizer that return a non-nil value will be wrapped with the
// function to create an InitilizedFunc.
type Initilizers []Initilizer

type Injectors []Injector

type Callback func([]reflect.Value)

func (injs Injectors) Inject(args []reflect.Value) (callbacks []Callback, err error) {
	var m lerr.Many
	out := slice.TransformSlice(injs, nil, func(i Injector, idx int) (Callback, bool) {
		cb, err := i.Inject(args)
		m.Add(err)
		return cb, cb != nil
	})

	return out, m.Cast()
}

// Apply the FuncInitilizers to fn. If fn is not a valid function, nil is
// returned.
func (fis Initilizers) Apply(fn any) *InitilizedFunc {
	v := reflect.ValueOf(fn)
	ft, err := newFunc(v)
	if err != nil {
		return nil
	}

	dis := slice.TransformSlice(fis, nil, func(i Initilizer, idx int) (Injector, bool) {
		out := i.Initilize(ft)
		return out, out != nil
	})

	return &InitilizedFunc{
		fnTypes: *ft,
		dis:     Injectors(dis),
		n:       ft.t.NumIn(),
		fn:      v,
	}
}

var (
	isInjectorFunc = filter.NumIn(filter.GT(0)).
		And(ltype.IsPtrToStruct.In(-1)).Filter
)

func newFunc(v reflect.Value) (*fnTypes, error) {
	t := v.Type()
	if !isInjectorFunc(t) {
		return nil, lerr.Str("not injector func")
	}

	dataType := t.In(t.NumIn() - 1).Elem()
	return &fnTypes{
		target: dataType,
		t:      t,
	}, nil
}

// InitilizedFunc wraps a function so that a set of FuncInjector will run before
// calling the function and any callbacks will run after it returns.
type InitilizedFunc struct {
	fnTypes
	dis Injectors
	n   int
	fn  reflect.Value
}

// Call the InitilizedFunc, running the each FuncInjector before the wrapped
// function is called and then invoking the callbacks when the wrapped function
// returns.
func (ifn *InitilizedFunc) Call(args []reflect.Value) []reflect.Value {
	if len(args) == ifn.n-1 {
		args = append(args, reflect.New(ifn.target))
	}

	// TODO: handle error
	cbs, _ := ifn.dis.Inject(args)

	out := ifn.fn.Call(args)

	for i := len(cbs) - 1; i >= 0; i-- {
		cbs[i](out)
	}

	return out
}

// Interface uses reflection to create a function without the target argument.
// This can be type cast to the function type and invoked as a normal function.
func (ifn *InitilizedFunc) Interface() any {
	in := make([]reflect.Type, ifn.n-1)
	for i := range in {
		in[i] = ifn.t.In(i)
	}
	out := make([]reflect.Type, ifn.t.NumOut())
	for i := range out {
		out[i] = ifn.t.Out(i)
	}
	t := reflect.FuncOf(in, out, false)
	return reflect.MakeFunc(t, ifn.Call).Interface()
}
