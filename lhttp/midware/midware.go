package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/util/linject"
)

type Midware struct {
	linject.FuncInitilizers
}

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
