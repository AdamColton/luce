package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lhttp"
)

type Decoder struct {
	lhttp.RequestDecoder
	FieldName string
}

func NewDecoder(d lhttp.RequestDecoder, fieldName string) Decoder {
	return Decoder{
		RequestDecoder: d,
		FieldName:      fieldName,
	}
}

type decoderInserter struct {
	lhttp.RequestDecoder
	idx []int
	t   reflect.Type
}

func (d Decoder) Initilize(t reflect.Type) DataInserter {
	if d.FieldName == "" {
		panic("Decoder.FieldName cannot be blank when used as Initilizer")
	}
	decField, hasDec := t.FieldByName(d.FieldName)
	if !hasDec {
		return nil
	}
	di := &decoderInserter{
		RequestDecoder: d.RequestDecoder,
		idx:            decField.Index,
	}

	di.t = decField.Type
	if di.t.Kind() != reflect.Ptr {
		panic("Decoder field should be pointer to struct:" + di.t.String())
	}
	di.t = di.t.Elem()
	if di.t.Kind() != reflect.Struct {
		panic("Decoder field should be pointer to struct")
	}
	return di
}

func (di *decoderInserter) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error) {
	v := reflect.New(di.t)
	err := di.Decode(v.Interface(), r)
	if err != nil {
		return nil, err
	}
	dst.FieldByIndex(di.idx).Set(v)
	return nil, nil
}
