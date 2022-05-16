package midware

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// URLFieldSetter references a Var in gorilla/mux.Vars. It fulfills
// FieldSetterInitilizer and FieldSetter. It is used to extract values from a
// URL and set them in the magic data field.
type URLFieldSetter struct {
	Var string
}

// Initilize fulfills FieldSetterInitilizer.
func (u URLFieldSetter) Initilize(fieldType reflect.Type) FieldSetter {
	return u
}

// Vars references mux.Vars. It is left exposed for testing.
var Vars = mux.Vars

// Set fulfills FieldSetter by setting field to a value from the URL using
// mux.Vars.
func (u URLFieldSetter) Set(w http.ResponseWriter, r *http.Request, field reflect.Value) (func(), error) {
	field.Set(reflect.ValueOf(Vars(r)[u.Var]))
	return nil, nil
}

// URL creates a magic Initilizer to extract a segment from the URL by name and
// set it to a magic data field by fieldName.
func URL(segment, fieldName string) Initilizer {
	return NewFieldInitilizer(URLFieldSetter{segment}, fieldName)
}
