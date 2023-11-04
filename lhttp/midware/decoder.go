package midware

import (
	"reflect"

	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

// NewDecoder creates a magic midware that decodes data from the request and
// sets it to a field on the magic data type. The field must be a pointer to a
// struct.
func NewDecoder(d lhttp.RequestDecoder, fieldName string) linject.FieldInitilizer {
	return NewFieldInitilizer(DecoderInitilizer{
		RequestDecoder: d,
	}, fieldName)
}

// DecoderInitilizer fulfills FieldSetterInitilizer.
type DecoderInitilizer struct {
	lhttp.RequestDecoder
}

var (
	decoderCheck = ltype.IsPtrToStruct.Check(filter.TypeErr("expected pointer to struct, got: %s"))
)

// Initilize fulfills FieldSetterInitilizer. It validates that the Type t is a
// pointer to a struct.
func (di DecoderInitilizer) InitilizeField(fn linject.Func, t reflect.Type) linject.FieldSetter {
	return &decoderSetter{
		RequestDecoder: di.RequestDecoder,
		Type:           decoderCheck.Panic(t).Elem(),
	}
}

type decoderSetter struct {
	lhttp.RequestDecoder
	reflect.Type
}

// Set fulfills FieldSetter. It creates and instance of the field to set, which
// will be a pointer to struct and calls Decode on the underlying RequestDecoder
// to set the field value.
func (ds decoderSetter) Set(args []reflect.Value, field reflect.Value) (func(), error) {
	_, r := GetWR(args)
	v := reflect.New(ds.Type)
	err := ds.Decode(v.Interface(), r)
	if err == nil {
		field.Set(v)
	}

	return nil, err
}
