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

// Midware wraps s as a linject.Field under s.FieldName, so a midware handler
// can take a *Session (or *lusers.User, via a pointer field of that type)
// as part of its data struct. It panics if s.FieldName is blank.
func (s *Store) Midware() linject.Field {
	if s.FieldName == "" {
		panic("Store.FieldName cannot be blank")
	}
	return midware.NewField(s, s.FieldName)
}

// Inject fulfills linject.FieldInitilizer by building the request's Session
// and arranging for it to be saved once the handler returns.
func (s *Store) Inject(w http.ResponseWriter, r *http.Request) (v any, fn func([]reflect.Value), err error) {
	var ses *Session
	ses, err = s.Session(w, r)
	if err == nil {
		v = ses
		fn = func(_ []reflect.Value) {
			err := ses.Save()
			if err != nil {
				// == projects.Code.luce.server ==
				// [ ] wire up real logging
				//  Inject's save-failure handler just prints to stdout.
				//  Replace once the server has a logging story.
				fmt.Println(err)
			}
		}
	}
	return
}

var sessionCheck = filter.IsType(reflector.Type[*Session]()).
	Check(filter.TypeErr("expected *lusess.Session, got: %s"))

// InitilizeField fulfills linject.FieldInitilizer. It only accepts a field
// of type *Session; it panics otherwise.
func (s *Store) InitilizeField(ft linject.FuncType, t reflect.Type) linject.FieldInjector {
	lerr.Must(sessionCheck(t))
	return linject.NewFieldInjector(s.Inject)
}
