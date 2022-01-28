package midware

import (
	"net/http"
	"reflect"
)

// Redirect is solves a side effect of the magic. When http.Redirect is invoked
// it prevents cookies from being written (and possibly other headers). Putting
// Redirect at the bottom of midware.NewMagic allows redirects to be set and
// invoked at the end of the call.
type Redirect struct {
	FieldName string
}

func (r Redirect) Initilize(t reflect.Type) DataInserter {
	if r.FieldName == "" {
		panic("Redirect.FieldName cannot be blank")
	}
	redirectField, hasRedirect := t.FieldByName(r.FieldName)
	if !hasRedirect {
		return nil
	}
	if redirectField.Type.Kind() != reflect.String {
		panic("Redirect field must be a string")
	}
	return &redirect{
		idx: redirectField.Index,
	}
}

type redirect struct {
	idx []int
}

func (to *redirect) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func(), error) {
	return func() {
		url := dst.Elem().FieldByIndex(to.idx).Interface().(string)
		if url != "" {
			http.Redirect(w, r, url, 302)
		}
	}, nil
}
