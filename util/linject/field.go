package linject

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
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
