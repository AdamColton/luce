package midware

import (
	"net/http"
	"reflect"
)

type Magic struct {
	Initilizers []Initilizer
}

func NewMagic(initilizers ...Initilizer) *Magic {
	return &Magic{
		Initilizers: initilizers,
	}
}

type Initilizer interface {
	Initilize(reflect.Type) DataInserter
}

type DataInserter interface {
	Insert(dst reflect.Value, r *http.Request) error
}

func (m *Magic) Handle(fn interface{}) http.HandlerFunc {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func || t.NumIn() != 3 {
		panic("Decode.Handler requires a func with 3 args")
	}

	dstType := t.In(2)
	useElem := true
	if dstType.Kind() == reflect.Ptr {
		dstType = dstType.Elem()
		useElem = false
	}

	var dis []DataInserter
	for _, i := range m.Initilizers {
		di := i.Initilize(dstType)
		if di != nil {
			dis = append(dis, di)
		}
	}

	v := reflect.ValueOf(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		dst := reflect.New(dstType)
		for _, di := range dis {
			di.Insert(dst, r)
		}

		if useElem {
			dst = dst.Elem()
		}
		v.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			dst,
		})
	}
}
