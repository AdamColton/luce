package midware

import (
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type Url struct {
	Var, FieldName string
}

type urlInserter struct {
	idx []int
	v   string
}

func (u Url) Initilize(t reflect.Type) DataInserter {
	if u.FieldName == "" {
		panic("Url.FieldName cannot be blank")
	}
	decField, hasDec := t.FieldByName(u.FieldName)
	if !hasDec {
		return nil
	}
	return &urlInserter{
		idx: decField.Index,
		v:   u.Var,
	}
}

func (ui *urlInserter) Insert(dst reflect.Value, r *http.Request) error {
	v := reflect.ValueOf(mux.Vars(r)[ui.v])
	dst.Elem().FieldByIndex(ui.idx).Set(v)
	return nil
}
