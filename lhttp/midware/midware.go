package midware

import (
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

func (m *Midware) Initilizer(i Initilizer) {
	m.FuncInitilizers = append(m.FuncInitilizers, wrappedInitilizer{i})
}

func (m *Midware) Handle(fn any) http.HandlerFunc {
	t := reflect.TypeOf(fn)
	if !linject.IsHttpHandler(t) {
		panic("Decode.Handler requires a func with 3 args")
	}

	ifn := m.Apply(fn)

	return ifn.Interface().(func(http.ResponseWriter, *http.Request))
}
