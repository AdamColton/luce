package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
)

// Midware holds the Initilizers that build the DataInserters for dataType.
type Midware struct {
	Initilizers []Initilizer
}

// New creates a set of midware initilizers that can be used to convert
// midwareFuncs to http.HandlerFuncs.
func New(initilizers ...Initilizer) *Midware {
	return &Midware{
		Initilizers: initilizers,
	}
}

// Initilizer takes a dataType and produces a DataInsert based on that type.
// The dataType will be taken from the third argument of a midwareFunc.
type Initilizer interface {
	Initilize(dataType reflect.Type) DataInserter
}

// DataInserter uses the ResponseWriter and Request to insert values into data.
type DataInserter interface {
	Insert(w http.ResponseWriter, r *http.Request, data reflect.Value) (callback func(), err error)
}

const (
	// ErrMidwareHandler is the panic value in Midware.Handle if midwareFunc
	// in not valid.
	ErrMidwareHandler = lerr.Str("Invalid MidwareFunc")
)

var (
	isStruct = filter.IsKind(reflect.Struct)
	isPtr    = filter.
			IsKind(reflect.Ptr)

	isPtrToStruct = isPtr.
			And(
			filter.Elem(
				isStruct,
			),
		)

	isResponseWriter = filter.IsNilRef((*http.ResponseWriter)(nil))
	isRequest        = filter.IsType((*http.Request)(nil))
	isMidwareFunc    = filter.NumIn(filter.EQ(3)).
				And(
			filter.In(0, isResponseWriter),
		).
		And(
			filter.In(1, isRequest),
		).And(
		filter.In(1, isPtrToStruct.Or(isStruct)),
	)

	midwareFuncCheck = filter.TypeCheck(isMidwareFunc, ErrMidwareHandler)
)

// Handle converts a midwareFunc to an http.HandlerFunc. The midwareFunc must be
// of the form
// - func(w http.ResponseWriter, r *http.Request, data struct{...})
func (m *Midware) Handle(midwareFunc interface{}) http.HandlerFunc {
	t := midwareFuncCheck.Panic(midwareFunc)

	dstType := t.In(2)
	ptr := isPtr(dstType)
	if ptr {
		dstType = dstType.Elem()
	}

	var dis []DataInserter
	for _, i := range m.Initilizers {
		di := i.Initilize(dstType)
		if di != nil {
			dis = append(dis, di)
		}
	}

	vfn := reflect.ValueOf(midwareFunc)

	return func(w http.ResponseWriter, r *http.Request) {
		dst := reflect.New(dstType)
		var callbacks []func()
		for _, di := range dis {
			callback, err := di.Insert(w, r, dst)
			lerr.Log(err)
			if callback != nil {
				callbacks = append(callbacks, callback)
			}
		}

		if !ptr {
			dst = dst.Elem()
		}
		vfn.Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
			dst,
		})

		for _, callback := range callbacks {
			callback()
		}
	}
}
