package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
)

// FieldSetterInitilizer is called by FieldInitilizer to handle any Type
// Reflection necessary in creating a FieldSetter. If no Type Reflection is
// necessary, a single type can fulfill FieldSetter, then FieldSetterInitilizer
// can just return itself.
type FieldSetterInitilizer interface {
	Initilize(fieldType reflect.Type) FieldSetter
}

// FieldSetter will set the value on 'set' generally from the request. This
// abstracts away the process of extracting a field from the magic data type.
type FieldSetter interface {
	Set(w http.ResponseWriter, r *http.Request, field reflect.Value) (callback func(), err error)
}

// FieldInitilizer fulfills Initilizer. It contains the logic for validating
// the FieldName exists and getting that field from the magic data type. This
// simplifies the process of creating midware that sets a value on one field.
type FieldInitilizer struct {
	FieldName string
	FieldSetterInitilizer
}

// NewFieldInitilizer with the given values.
func NewFieldInitilizer(fi FieldSetterInitilizer, fieldName string) FieldInitilizer {
	return FieldInitilizer{
		FieldName:             fieldName,
		FieldSetterInitilizer: fi,
	}
}

type fieldInserter struct {
	idx []int
	FieldSetter
}

const (
	// ErrFieldName is the panic value used by FieldInitilizer.Initilize if the
	// FieldName is not valid.
	ErrFieldName = lerr.Str("Invalid FieldSetter.FieldName")
)

var (
	fieldName      = filter.MustRegex(`\p{Lu}(\p{L}|\p{N}|_)*`)
	checkFieldName = fieldName.Check(ErrFieldName)
)

// Initilize fulfills Initilizer and validates the FieldName is not blank and
// checks if magic data type 't' has the field. If the field exists a
// DataInserter is created that will invoke the FieldSetter on that field.
func (fs FieldInitilizer) Initilize(dataType reflect.Type) DataInserter {
	checkFieldName.Panic(fs.FieldName)
	field, hasfield := dataType.FieldByName(fs.FieldName)
	if !hasfield {
		return nil
	}
	return &fieldInserter{
		idx:         field.Index,
		FieldSetter: fs.FieldSetterInitilizer.Initilize(field.Type),
	}
}

// Insert fulfills DataInserter. It gets the field from the magic data value and
// passes it into the field setter.
func (fi *fieldInserter) Insert(w http.ResponseWriter, r *http.Request, data reflect.Value) (func(), error) {
	return fi.Set(w, r, data.Elem().FieldByIndex(fi.idx))
}
