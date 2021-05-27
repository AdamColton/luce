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
	urlField, hasUrl := t.FieldByName(u.FieldName)
	if !hasUrl {
		return nil
	}
	return &urlInserter{
		idx: urlField.Index,
		v:   u.Var,
	}
}

func (ui *urlInserter) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) error {
	v := reflect.ValueOf(mux.Vars(r)[ui.v])
	dst.Elem().FieldByIndex(ui.idx).Set(v)
	return nil
}
