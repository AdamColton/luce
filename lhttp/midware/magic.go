package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

// Magic allows for a third "magic" argument to be passed to http handlers.
type Magic struct {
	Initilizers []Initilizer
}

func NewMagic(initilizers ...Initilizer) *Magic {
	return &Magic{
		Initilizers: initilizers,
	}
}

// Initilizer only runs when a route is added, not on each invocation. If the
// route meets the parameters, the returned DataInserter will be invoked before
// the route is called.
type Initilizer interface {
	Initilize(reflect.Type) DataInserter
}

// DataInserter will typically modify the dst value. It cann also return a
// callback that will run after the handler is closed.
type DataInserter interface {
	Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error)
}

// Handle requires fn to be a function with 3 arguments. The frist should be
// http.ResponseWriter, the second *http.Request and the third should be a
// pointer to a struct. The struct will be passed into the initilizers to setup
// the DataInserters.
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
		callbacks := make([]func(), 0, len(dis))
		for _, di := range dis {
			callback, err := di.Insert(w, r, dst)
			lerr.Log(err)
			if callback != nil {
				callbacks = append(callbacks, callback)
			}
		}

		if useElem {
			dst = dst.Elem()
		}
		v.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			dst,
		})

		for _, callback := range callbacks {
			callback()
		}
	}
}
