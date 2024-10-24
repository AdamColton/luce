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
	Counters  lmap.Wrapper[string, int]
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
	ce.Counters = lmap.FromIter(ce.Endpoints.Keys(nil).Iter(), func(str string, idx int) (string, int, bool) {
		return str, 0, true
	})
	return
}

var zeroCounter filter.MapValueFilter[string, int] = func(i int) bool {
	return i == 0
}

// UnusedEnpoints panics if any endpoints are unused
func (ce ClientEndpoints) UnusedEnpoints() {
	unused := zeroCounter.KeySlice(ce.Counters.Map())
	if len(unused) > 0 {
		panic("Unused Handlers: " + strings.Join(unused, "; "))
	}
}

// Add pops the endpoints associated with the handler and adds it to the client
// with the given route. Ths allows UnusedEnpoints to be called which will
// panic if any handlers are unused.
func (ce ClientEndpoints) Add(handler string, rCfg *RouteConfig) {
	rr := ce.Endpoints.GetVal(handler)
	if rr == nil {
		panic("coud not find handler " + handler)
	}
	ce.Counters.Set(handler, ce.Counters.GetVal(handler)+1)
	ce.Client.Add(rr, rCfg)
}
