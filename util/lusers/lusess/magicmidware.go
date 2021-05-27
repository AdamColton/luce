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

func (s *Store) Initilize(t reflect.Type) midware.Injector {
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

func (mi midwareInserter) Inject(w http.ResponseWriter, r *http.Request, dst reflect.Value) (func([]reflect.Value), error) {
	s, err := mi.s.Session(w, r)
	if err != nil {
		return nil, err
	}
	dst.Elem().FieldByIndex(mi.idx).Set(reflect.ValueOf(s))
	return func(rets []reflect.Value) {
		s.Save()
	}, nil
}
