package linject

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// FieldInitilizer is called by FieldInitilizer to handle any Type
// Reflection necessary in creating a FieldSetter. If no Type Reflection is
// necessary, a single type can fulfill FieldSetter, then FieldSetterInitilizer
// can just return itself.
type FieldInitilizer interface {
	InitilizeField(FuncType, reflect.Type) FieldInjector
}

// FieldInjector will set the value on 'set' generally from the request. This
// abstracts away the process of extracting a field from the magic data type.
type FieldInjector interface {
	InjectField(args []reflect.Value, field reflect.Value) (callback func([]reflect.Value), err error)
}

// FieldInitilizer fulfills Initilizer. It contains the logic for validating
// the FieldName exists and getting that field from the magic data type. This
// simplifies the process of creating midware that sets a value on one field.
type Field struct {
	FieldName string
	FieldInitilizer
	FuncType, FieldType filter.Type
}

// NewField with the given values.
func NewField(fi FieldInitilizer, fieldName string) Field {
	return Field{
		FieldName:       fieldName,
		FieldInitilizer: fi,
	}
}

type fieldInserter struct {
	argsIdx int
	idx     []int
	FieldInjector
}

var (
	FieldName      = lerr.Must(filter.Regex(`\p{Lu}(\p{L}|\p{N}|_)*`))
	checkFieldName = FieldName.Check(func(s string) error {
		return lerr.Str("Invalid FieldName: " + s)
	})
)

// Initilize fulfills Initilizer and validates the FieldName is not blank and
// checks if magic data type 't' has the field. If the field exists a
// DataInserter is created that will invoke the FieldSetter on that field.
func (fs Field) Initilize(fn FuncType) Injector {
	in := fn.Fn().NumIn()
	checkFieldName.Panic(fs.FieldName)
	field, hasfield := fn.Target().FieldByName(fs.FieldName)
	fit, fnt := fs.FieldType.Filter, fs.FuncType.Filter
	if !hasfield || (fnt != nil && !fnt(fn.Fn())) || (fit != nil && !fit(field.Type)) {
		return nil
	}
	setter := fs.FieldInitilizer.InitilizeField(fn, field.Type)
	if setter == nil {
		return nil
	}
	return &fieldInserter{
		argsIdx:       in - 1,
		idx:           field.Index,
		FieldInjector: setter,
	}
}

func (fi *fieldInserter) Inject(args []reflect.Value) (callback func([]reflect.Value), err error) {
	return fi.InjectField(args, args[fi.argsIdx].Elem().FieldByIndex(fi.idx))
}

type fieldSetter struct {
	set func(args []reflect.Value, field reflect.Value) (callback func([]reflect.Value), err error)
}

func (ofs fieldSetter) InjectField(args []reflect.Value, field reflect.Value) (callback func([]reflect.Value), err error) {
	return ofs.set(args, field)
}

var (
	fnType           = reflector.Type[func([]reflect.Value)]()
	fieldSetterCheck = filter.NumOutEq(3).
				And(filter.IsType(fnType).Out(1)).
				And(filter.IsType(ltype.Err).Out(2)).
				Check(filter.TypeErr("NewFieldInjector expected func([]reflect.Value)(t T,callback func([]reflect.Value), err error), got %s"))

	valueSliceType = reflector.Type[[]reflect.Value]()
	wrappedArgsFn  = filter.NumInEq(1).
			And(filter.IsType(valueSliceType).In(0)).Filter
)

// NewFieldInjector takes a function and converts it to a field setter. The
// function must have 3 returns. The first is the value the field will be set
// to. The second is the callback function and the third is an error.
//
// The arguments can either have a single argument of []reflect.Value in which
// case the arguments will be passed along. Or if it expects a specific argument
// pattern for the function, it can match the leading arguments. For instance, a
// FieldSetter on an HttpHandler could have arguments of (w http.ResponseWriter,
// r *http.Request)
func NewFieldInjector(fn any) FieldInjector {
	t := lerr.Must(fieldSetterCheck(fn))
	fnv := reflect.ValueOf(fn)
	var getArgs func([]reflect.Value) []reflect.Value
	if wrappedArgsFn(t) {
		getArgs = func(args []reflect.Value) []reflect.Value {
			return []reflect.Value{reflect.ValueOf(args)}
		}
	} else {
		getArgs = func(args []reflect.Value) []reflect.Value {
			return args[:t.NumIn()]
		}
	}
	return fieldSetter{
		set: func(args []reflect.Value, field reflect.Value) (callback func([]reflect.Value), err error) {
			fnArgs := getArgs(args)
			out := fnv.Call(fnArgs)
			reflector.Set(field, out[0])
			callback = out[1].Interface().(func([]reflect.Value))
			i := out[2].Interface()
			if i != nil {
				err = i.(error)
			}
			return
		},
	}
}

type setterWrapper struct {
	setter any
}

func (fs setterWrapper) InitilizeField(ft FuncType, t reflect.Type) FieldInjector {
	return NewFieldInjector(fs.setter)
}

func NewFieldSetter(setter any, fieldName string, funcType, fieldType filter.Type) Field {
	fi := NewField(setterWrapper{setter}, fieldName)
	fi.FuncType = funcType
	fi.FieldType = fieldType
	return fi
}
