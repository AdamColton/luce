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

func (m *Midware) Inits(initilizers ...Initilizer) *Midware {
	for _, i := range initilizers {
		m.Initilizers = append(m.Initilizers, wrappedInitilizer{i})
	}
	return m
}

func (m *Midware) Handle(fn any) http.HandlerFunc {
	t := reflect.TypeOf(fn)
	if !HttpHandlerType.Filter(t) {
		panic(fmt.Errorf("invalid Midware funce: %s", t))
	}

	ifn := m.Apply(fn)

	return ifn.Interface().(func(http.ResponseWriter, *http.Request))
}
