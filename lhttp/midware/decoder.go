package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lhttp"
)

type Decoder struct {
	lhttp.RequestDecoder
	lhttp.ErrHandler
	FieldName string
}

func NewDecoder(d lhttp.RequestDecoder, fieldName string) Decoder {
	return Decoder{
		RequestDecoder: d,
		FieldName:      fieldName,
	}
}

// Handler expects fn to be a function with 3 arguments. The first should be
// http.ResponseWriter and the second should be *http.Request.  The third should
// be the type that will be populated by the form.
func (d Decoder) Handler(fn interface{}) http.HandlerFunc {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func || t.NumIn() != 3 {
		panic("Decode.Handler requires a func with 3 args")
	}

	dstType := t.In(2)
	if dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
	}
	v := reflect.ValueOf(fn)
	return func(w http.ResponseWriter, r *http.Request) {
		dst := reflect.New(dstType)
		err := d.Decode(dst.Interface(), r)
		if d.Check(w, r, err) {
			return
		}
		v.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			dst,
		})
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

func (di *decoderInserter) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) error {
	v := reflect.New(di.t)
	err := di.Decode(v.Interface(), r)
	if err != nil {
		return err
	}
	dst.Elem().FieldByIndex(di.idx).Set(v)
	return nil
}
