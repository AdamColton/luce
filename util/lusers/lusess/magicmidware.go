package lusess

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lhttp/midware"
)

type midwareInserter struct {
	s   *Store
	idx []int
}

func (s *Store) Initilize(t reflect.Type) midware.DataInserter {
	if s.FieldName == "" {
		panic("Store.FieldName cannot be blank")
	}
	sField, has := t.FieldByName(s.FieldName)
	if !has {
		return nil
	}
	return &midwareInserter{
		s:   s,
		idx: sField.Index,
	}
}

func (mi midwareInserter) Insert(w http.ResponseWriter, r *http.Request, dst reflect.Value) (error, func()) {
	s, err := mi.s.Session(w, r)
	if err != nil {
		return err, nil
	}
	dst.Elem().FieldByIndex(mi.idx).Set(reflect.ValueOf(s))
	return nil, func() {
		s.Save()
	}
}
