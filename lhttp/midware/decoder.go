package midware

import (
	"net/http"
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
	ds := &decoderSetter{
		RequestDecoder: di.RequestDecoder,
		Type:           decoderCheck.Panic(t).Elem(),
	}
	return linject.NewFieldSetter(ds.set)
}

type decoderSetter struct {
	lhttp.RequestDecoder
	reflect.Type
}

func (ds decoderSetter) set(w http.ResponseWriter, r *http.Request) (any, func(), error) {
	v := reflect.New(ds.Type)
	err := ds.Decode(v.Interface(), r)

	return v.Interface(), nil, err
}
