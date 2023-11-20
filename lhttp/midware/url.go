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
var Vars = mux.Vars

// URL creates a magic Initilizer to extract a segment from the URL by name and
// set it to a magic data field by fieldName.
func URL(segment, fieldName string) linject.FieldInitilizer {
	fi := NewFieldInitilizer(URLFieldSetter{segment}, fieldName)
	fi.FieldType = filter.IsKind(reflect.String)
	return fi
}

// Initilize fulfills FieldSetterInitilizer.
func (u URLFieldSetter) InitilizeField(fn linject.Func, t reflect.Type) linject.FieldSetter {
	return linject.NewFieldSetter(u.set)
}

func (u URLFieldSetter) set(w http.ResponseWriter, r *http.Request) (string, func(), error) {
	return Vars(r)[u.Var], nil, nil
}
