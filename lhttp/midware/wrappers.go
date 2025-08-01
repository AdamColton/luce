package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/linject"
)

type Initilizer interface {
	Initilize(reflect.Type) Injector
}

type Injector interface {
	Inject(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func([]reflect.Value), error)
}

type wrappedInitilizer struct {
	Initilizer
}

func (wi wrappedInitilizer) Initilize(fn linject.FuncType) linject.Injector {
	di := wi.Initilizer.Initilize(fn.Target())
	if di == nil {
		return nil
	}
	return wrappedInjector{di}
}

type wrappedInjector struct {
	Injector
}

func (wdi wrappedInjector) Inject(args []reflect.Value) (callback func([]reflect.Value), err error) {
	w := args[0].Interface().(http.ResponseWriter)
	r := args[1].Interface().(*http.Request)
	d := args[2]
	if wdi.Injector == nil {
		return nil, lerr.Str("nil injector")
	}
	return wdi.Injector.Inject(w, r, d)
}
