package midware

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/linject"
)

// Midware holds the Initilizers that build the DataInserters for dataType.
type Midware struct {
	linject.FuncInitilizers
}

// New creates a set of midware initilizers that can be used to convert
// midwareFuncs to http.HandlerFuncs.
func New(initilizers ...linject.FuncInitilizer) *Midware {
	return &Midware{
		FuncInitilizers: initilizers,
	}
}

func (m *Midware) Initilizers(initilizers ...Initilizer) *Midware {
	for _, i := range initilizers {
		m.FuncInitilizers = append(m.FuncInitilizers, wrappedInitilizer{i})
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
