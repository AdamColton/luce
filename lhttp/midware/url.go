package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/gorilla/mux"
)

// URLFieldSetter references a Var in gorilla/mux.Vars. It fulfills
// FieldSetterInitilizer and FieldSetter. It is used to extract values from a
// URL and set them in the magic data field.
type URLFieldSetter struct {
	Var string
}

// Vars references mux.Vars. It is left exposed for testing.
var (
	Vars     = mux.Vars
	isString = filter.IsKind(reflect.String)
)

// URL creates a magic Initilizer to extract a segment from the URL by name and
// set it to a magic data field by fieldName.
func URL(segment, fieldName string) linject.Field {
	fn := func(w http.ResponseWriter, r *http.Request) (string, func([]reflect.Value), error) {
		return Vars(r)[segment], nil, nil
	}
	return linject.NewFieldSetter(fn, fieldName, filter.AnyType(), isString)
}
