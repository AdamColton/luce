package service

import (
	"strings"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
)

var RequestResponderFilter = filter.ConvertableTo[RequestResponder]()

func epMethodToKV(m *reflector.Method, idx int) (string, RequestResponder, bool) {
	fn := m.Func.Interface().(func(*Request) *Response)
	return m.Name, fn, true
}

type ClientEndpoints struct {
	*Client
	Endpoints lmap.Wrapper[string, RequestResponder]
}

// NewClientEndpoints solves a specific problem. Creating an endpoint function
// but forgetting to register it with the Client. ClientEndpoints finds all
// methods on the endpoints object that fulfill RequestResponder. As each is
// used to add a route to the Client it is removed the Endpoints map. Calling
// UnusedEnpoints will panic if a RequestResponder was added but not used in a
// route.
func NewClientEndpoints(addr string, endpoints any) (ce ClientEndpoints, err error) {
	ce.Client, err = NewClient(addr)
	if err != nil {
		return
	}
	epMethods, _ := RequestResponderFilter.Method().SliceInPlace(reflector.MethodsOn(endpoints))
	ce.Endpoints = lmap.FromIter(epMethods.Iter(), epMethodToKV)
	return
}

// UnusedEnpoints panics if any endpoints are unused
func (ce ClientEndpoints) UnusedEnpoints() {
	if ce.Endpoints.Len() != 0 {
		panic("Unused Handlers: " + strings.Join(ce.Endpoints.Keys(nil), "; "))
	}
}

// Add pops the endpoints associated with the handler and adds it to the client
// with the given route. Ths allows UnusedEnpoints to be called which will
// panic if any handlers are unused.
func (ce ClientEndpoints) Add(handler string, r *RouteConfig) {
	ce.Client.Add(ce.Endpoints.MustPop(handler), r)
}
