package lusess

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/linject"
	"github.com/adamcolton/luce/util/reflector"
)

func (s *Store) Midware() linject.Field {
	if s.FieldName == "" {
		panic("Store.FieldName cannot be blank")
	}
	return midware.NewField(s, s.FieldName)
}

func (s *Store) Inject(w http.ResponseWriter, r *http.Request) (v any, fn func([]reflect.Value), err error) {
	var ses *Session
	ses, err = s.Session(w, r)
	if err == nil {
		v = ses
		fn = func(_ []reflect.Value) {
			err := ses.Save()
			if err != nil {
				// TODO: use logging
				fmt.Println(err)
			}
		}
	}
	return
}

var sessionCheck = filter.IsType(reflector.Type[*Session]()).
	Check(filter.TypeErr("expected *lusess.Session, got: %s"))

func (s *Store) InitilizeField(ft linject.FuncType, t reflect.Type) linject.FieldInjector {
	lerr.Must(sessionCheck(t))
	return linject.NewFieldInjector(s.Inject)
}
