package midware

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/linject"
)

// Midware holds the Initilizers that build the DataInserters for dataType.
type Midware struct {
	linject.Initilizers
}

// New creates a set of midware initilizers that can be used to convert
// midwareFuncs to http.HandlerFuncs.
func New(initilizers ...linject.Initilizer) *Midware {
	return &Midware{
		Initilizers: initilizers,
	}
}

// Inits adds initilizers, each wrapped to fulfill linject.Initilizer, and
// returns m for chaining.
func (m *Midware) Inits(initilizers ...Initilizer) *Midware {
	for _, i := range initilizers {
		m.Initilizers = append(m.Initilizers, wrappedInitilizer{i})
	}
	return m
}

// Handle converts fn, a MidwareFunc of the form
// func(w http.ResponseWriter, r *http.Request, data *struct{...}), to an
// http.HandlerFunc using m's Initilizers to populate data. It panics if fn
// does not have that shape.
func (m *Midware) Handle(fn any) http.HandlerFunc {
	t := reflect.TypeOf(fn)
	if !HttpHandlerType.Filter(t) {
		panic(fmt.Errorf("invalid Midware func: %s", t))
	}

	ifn := m.Apply(fn)

	return ifn.Interface().(func(http.ResponseWriter, *http.Request))
}
