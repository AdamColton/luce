package linject

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// FieldInitilizer is called by Field to handle any type reflection necessary in
// creating a FieldInjector. If no reflection is necessary, a single type can
// fulfill both FieldInitilizer and FieldInjector, and InitilizeField can just
// return itself.
type FieldInitilizer interface {
	InitilizeField(FuncType, reflect.Type) FieldInjector
}

// FieldInjector injects a value into one field of the target, generally from the
// request. It is given the arguments of the function and the field to set. This
// abstracts away the process of extracting a field from the target. The callback
// it returns, if any, is run with the return values of the function.
type FieldInjector interface {
	InjectField(args []reflect.Value, field reflect.Value) (callback func([]reflect.Value), err error)
}

// Field fulfills Initilizer. It contains the logic for validating that the
// FieldName exists on the target and getting that field from it. This simplifies
// the process of creating midware that sets a value on one field. FuncType and
// FieldType, when set, limit the functions and the fields it applies to.
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
	// FieldName is a filter for valid exported field names.
	FieldName      = lerr.Must(filter.Regex(`\p{Lu}(\p{L}|\p{N}|_)*`))
	checkFieldName = FieldName.Check(func(s string) error {
		return lerr.Str("Invalid FieldName: " + s)
	})
)

// Initilize fulfills Initilizer. It panics if FieldName is not a valid exported
// field name. It returns nil if the target has no such field, if FuncType or
// FieldType reject the function or the field, or if the FieldInitilizer returns
// nil. Otherwise the Injector it returns invokes the FieldInjector on the field.
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

// NewFieldSetter creates a Field for fieldName that calls setter, which must have
// the form NewFieldInjector takes. The funcType and fieldType filters limit the
// functions and fields it applies to; an empty filter.Type accepts everything.
func NewFieldSetter(setter any, fieldName string, funcType, fieldType filter.Type) Field {
	fi := NewField(setterWrapper{setter}, fieldName)
	fi.FuncType = funcType
	fi.FieldType = fieldType
	return fi
}
