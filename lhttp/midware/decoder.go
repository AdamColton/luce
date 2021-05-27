package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lhttp"
)

type Decoder struct {
	lhttp.RequestDecoder
	lhttp.ErrHandler
}

func NewDecoder(d lhttp.RequestDecoder) Decoder {
	return Decoder{
		RequestDecoder: d,
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
