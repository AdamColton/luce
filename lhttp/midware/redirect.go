package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/filter"
)

// Redirect is solves a side effect of the magic. When http.Redirect is invoked
// it prevents cookies from being written (and possibly other headers). Putting
// Redirect at the bottom of midware.NewMagic allows redirects to be set and
// invoked at the end of the call.
type Redirect struct{}

var checkRedirect = filter.TypeCheck(filter.IsKind(reflect.String), typeErr("Invalid Redirect field: "))

func NewRedirect(fieldName string) Initilizer {
	return NewFieldInitilizer(Redirect{}, fieldName)
}

// Initilize fulfills FieldSetterInitilizer.
func (rd Redirect) Initilize(fieldType reflect.Type) FieldSetter {
	checkRedirect.Panic(fieldType)
	return rd
}

// Set fulfills FieldSetter by setting field to a value from the URL using
// mux.Vars.
func (rd Redirect) Set(w http.ResponseWriter, r *http.Request, field reflect.Value) (func(), error) {
	return func() {
		url := field.Interface().(string)
		if url != "" {
			http.Redirect(w, r, url, 302)
		}
	}, nil
}
