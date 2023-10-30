package linject

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

type FuncInitilizer interface {
	Initilize(Func) DataInserter
}

type DataInserter interface {
	Insert(args []reflect.Value) (callback func(), err error)
}

type FuncInitilizers []FuncInitilizer

type Func interface {
	DataType() reflect.Type
	Fn() reflect.Type

	private()
}

type fnc struct {
	dataType reflect.Type
	t        reflect.Type
}

func (ft *fnc) DataType() reflect.Type {
	return ft.dataType
}

func (ft *fnc) Fn() reflect.Type {
	return ft.t
}

func (ft *fnc) private() {
}

var (
	isInjectorFunc = filter.NumIn(filter.GT(0)).
		And(ltype.IsPtrToStruct.In(-1)).Filter
)

func NewFunc(fn any) (Func, error) {
	return newFunc(reflect.ValueOf(fn))
}

func newFunc(v reflect.Value) (*fnc, error) {
	t := v.Type()
	if !isInjectorFunc(t) {
		return nil, lerr.Str("not injector func")
	}

	dataType := t.In(t.NumIn() - 1).Elem()
	return &fnc{
		dataType: dataType,
		t:        t,
	}, nil
}

func (fis FuncInitilizers) Apply(fn any) *InitilizedFunc {
	v := reflect.ValueOf(fn)
	ft, err := newFunc(v)
	if err != nil {
		return nil
	}
	var dis slice.Slice[DataInserter]
	for _, i := range fis {
		dis = dis.AppendNotZero(i.Initilize(ft))
	}

	return &InitilizedFunc{
		fnc: *ft,
		dis: dis,
		n:   ft.t.NumIn(),
		fn:  v,
	}
}

type InitilizedFunc struct {
	fnc
	dis []DataInserter
	n   int
	fn  reflect.Value
}

func (ifn *InitilizedFunc) Call(args []reflect.Value) []reflect.Value {
	if len(args) == ifn.n-1 {
		args = append(args, reflect.New(ifn.dataType))
	}

	var cbs slice.Slice[func()]
	for _, di := range ifn.dis {
		cb, _ := di.Insert(args)
		cbs = cbs.AppendNotZero(cb)
	}

	out := ifn.fn.Call(args)

	for i := len(cbs) - 1; i >= 0; i-- {
		cbs[i]()
	}

	return out
}

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
