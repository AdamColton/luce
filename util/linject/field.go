package linject

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// FieldSetterInitilizer is called by FieldInitilizer to handle any Type
// Reflection necessary in creating a FieldSetter. If no Type Reflection is
// necessary, a single type can fulfill FieldSetter, then FieldSetterInitilizer
// can just return itself.
type FieldSetterInitilizer interface {
	InitilizeField(Func, reflect.Type) FieldSetter
}

// FieldSetter will set the value on 'set' generally from the request. This
// abstracts away the process of extracting a field from the magic data type.
type FieldSetter interface {
	Set(args []reflect.Value, field reflect.Value) (callback func(), err error)
}

// FieldInitilizer fulfills Initilizer. It contains the logic for validating
// the FieldName exists and getting that field from the magic data type. This
// simplifies the process of creating midware that sets a value on one field.
type FieldInitilizer struct {
	FieldName string
	FieldSetterInitilizer
	FuncType, FieldType filter.Type
}

// NewFieldInitilizer with the given values.
func NewFieldInitilizer(fi FieldSetterInitilizer, fieldName string) FieldInitilizer {
	return FieldInitilizer{
		FieldName:             fieldName,
		FieldSetterInitilizer: fi,
	}
}

type fieldInserter struct {
	argsIdx int
	idx     []int
	FieldSetter
}

var (
	FieldName      = filter.MustRegex(`\p{Lu}(\p{L}|\p{N}|_)*`)
	checkFieldName = FieldName.Check(func(s string) error {
		return lerr.Str("Invalid FieldName: " + s)
	})
)

// Initilize fulfills Initilizer and validates the FieldName is not blank and
// checks if magic data type 't' has the field. If the field exists a
// DataInserter is created that will invoke the FieldSetter on that field.
func (fs FieldInitilizer) Initilize(fn Func) DataInserter {
	in := fn.Fn().NumIn()
	checkFieldName.Panic(fs.FieldName)
	field, hasfield := fn.DataType().FieldByName(fs.FieldName)
	fit, fnt := fs.FieldType.Filter, fs.FuncType.Filter
	if !hasfield || (fnt != nil && !fnt(fn.Fn())) || (fit != nil && !fit(field.Type)) {
		return nil
	}
	setter := fs.FieldSetterInitilizer.InitilizeField(fn, field.Type)
	if setter == nil {
		return nil
	}
	return &fieldInserter{
		argsIdx:     in - 1,
		idx:         field.Index,
		FieldSetter: setter,
	}
}

func (fi *fieldInserter) Insert(args []reflect.Value) (callback func(), err error) {
	return fi.Set(args, args[fi.argsIdx].Elem().FieldByIndex(fi.idx))
}

type outFieldSetter struct {
	set func(args []reflect.Value, field reflect.Value) (callback func(), err error)
}

func (ofs outFieldSetter) Set(args []reflect.Value, field reflect.Value) (callback func(), err error) {
	return ofs.set(args, field)
}

var (
	fnType              = reflector.Type[func()]()
	outFieldSetterCheck = filter.NumOutEq(3).
				And(filter.IsType(fnType).Out(1)).
				And(filter.IsType(ltype.Err).Out(2)).
				Check(filter.TypeErr("exptedted func([]reflect.Value)(t T,callback func(), err error), got %s"))

	valueSliceType = reflector.Type[[]reflect.Value]()
	wrappedArgsFn  = filter.NumInEq(1).
			And(filter.IsType(valueSliceType).In(0)).Filter
)

// NewFieldSetter takes a function and converts it to a field setter. The
// function must have 3 returns. The first is the value the field will be set
// to. The second is the callback function and the third is an error.
//
// The arguments can either have a single argument of []reflect.Value in which
// case the arguments will be passed along. Or if it expects a specific argument
// pattern for the function, it can match the leading arguments. For instance, a
// FieldSetter on an HttpHandler could have arguments of (w http.ResponseWriter,
// r *http.Request)
func NewFieldSetter(fn any) FieldSetter {
	outFieldSetterCheck.Panic(fn)
	fnv := reflect.ValueOf(fn)
	var getArgs func([]reflect.Value) []reflect.Value
	t := fnv.Type()
	if wrappedArgsFn(t) {
		getArgs = func(args []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(args)}
		}
	} else {
		getArgs = func(args []reflect.Value) []reflect.Value {
			return args[:t.NumIn()]
		}
	}
	return outFieldSetter{
		set: func(args []reflect.Value, field reflect.Value) (callback func(), err error) {
			fnArgs := getArgs(args)
			out := fnv.Call(fnArgs)
			reflector.Set(field, out[0])
			callback = out[1].Interface().(func())
			i := out[2].Interface()
			if i != nil {
				err = i.(error)
			}
			return
		},
	}
}
