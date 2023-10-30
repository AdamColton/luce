package linject

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

type Initilizer interface {
	Initilize(dataType reflect.Type) Injector
}

type Injector interface {
	Inject(data reflect.Value) (callback func(), err error)
}

type Initilizers []Initilizer

type Injectors struct {
	t    reflect.Type
	injs []Injector
}

func (inits Initilizers) Initilize(dataType reflect.Type) *Injectors {
	var injs slice.Slice[Injector]
	for _, i := range inits {
		injs = injs.AppendNotZero(i.Initilize(dataType))
	}
	if len(injs) == 0 {
		return nil
	}
	return &Injectors{
		injs: injs,
		t:    dataType,
	}
}

func (injs Injectors) Inject(data reflect.Value) (callbacks []func(), err error) {
	if data.Type() != injs.t {
		return nil, lerr.Str("types do not match")
	}
	var errs lerr.Many
	var cbs slice.Slice[func()]
	for _, i := range injs.injs {
		cb, err := i.Inject(data)
		errs = errs.Add(err)
		cbs = cbs.AppendNotZero(cb)
	}
	return cbs, errs.Cast()
}

func CheckField(on reflect.Type, name string, t reflect.Type) *reflect.StructField {
	on = ltype.CheckStructPtr(on)
	if on == nil {
		return nil
	}

	sf, ok := on.FieldByName(name)
	if !ok || sf.Type != t {
		return nil
	}
	return &sf
}
