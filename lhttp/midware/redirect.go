package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
)

// Redirect is solves a side effect of the callbacks. When http.Redirect is
// invoked it prevents cookies from being written (and possibly other headers).
// Putting Redirect at the bottom of midware.NewMagic allows redirects to be set
// and invoked at the end of the call.
type Redirect struct{}

var checkRedirect = filter.IsKind(reflect.String).
	Check(filter.TypeErr("expected string, got: %s"))

// NewRedirect creates a Redirect Initilizer. It should be the last Initilizer
// in the Midware.
func NewRedirect(fieldName string) linject.Field {
	return NewField(Redirect{}, fieldName)
}

// Initilize fulfills FieldSetterInitilizer.
func (rd Redirect) InitilizeField(fn linject.FuncType, t reflect.Type) linject.FieldInjector {
	lerr.Must(checkRedirect(t))
	return rd
}

// Set fulfills FieldSetter by setting field to a value from the URL using
// mux.Vars.
func (rd Redirect) InjectField(args []reflect.Value, field reflect.Value) (func([]reflect.Value), error) {
	w, r := GetWR(args)
	return func(rets []reflect.Value) {
		url := field.Interface().(string)
		if url != "" {
			http.Redirect(w, r, url, http.StatusFound)
		}
	}, nil
}
