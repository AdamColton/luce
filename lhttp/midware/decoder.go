package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/util/filter"
)

// NewDecoder creates a magic midware that decodes data from the request and
// sets it to a field on the magic data type.
func NewDecoder(d lhttp.RequestDecoder, fieldName string) Initilizer {
	return NewFieldInitilizer(DecoderInitilizer{
		RequestDecoder: d,
	}, fieldName)
}

// DecoderInitilizer fulfills FieldSetterInitilizer.
type DecoderInitilizer struct {
	lhttp.RequestDecoder
}

const (
	// ErrDecoderField is the panic value used by DecoderInitilizer if the field
	// type is not a pointer to a struct.
	ErrDecoderField = lerr.Str("Decoder field should be pointer to struct")
)

var decoderCheck = filter.TypeCheck(isPtrToStruct, ErrDecoderField)

// Initilize fulfills FieldSetterInitilizer. It validates that the Type t is a
// pointer to a struct.
func (di DecoderInitilizer) Initilize(fieldType reflect.Type) FieldSetter {
	return &decoderSetter{
		RequestDecoder: di.RequestDecoder,
		Type:           decoderCheck.Panic(fieldType).Elem(),
	}
}

type decoderSetter struct {
	lhttp.RequestDecoder
	reflect.Type
}

// Set fulfills FieldSetter. It creates and instance of the field to set, which
// will be a pointer to struct and calls Decode on the underlying RequestDecoder
// to set the field value.
func (ds decoderSetter) Set(w http.ResponseWriter, r *http.Request, field reflect.Value) (func(), error) {
	v := reflect.New(ds.Type)
	err := ds.Decode(v.Interface(), r)
	if err == nil {
		field.Set(v)
	}

	return nil, err
}
