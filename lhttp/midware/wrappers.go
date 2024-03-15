package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/linject"
)

// Initilizer builds an Injector for a given data type. Midware.Inits wraps
// each Initilizer to fulfill linject.Initilizer.
type Initilizer interface {
	Initilize(reflect.Type) Injector
}

// Injector sets a value on dst from the request, returning an optional
// callback to run after the handler completes. It mirrors linject.Injector,
// but with concrete http.ResponseWriter/*http.Request arguments in place of
// reflected ones.
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
	return wdi.Injector.Inject(w, r, d)
}
