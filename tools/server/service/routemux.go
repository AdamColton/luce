package service

type RequestHandler func(r Request)

type RequestResponder func(r Request) Response

type Mux struct {
	Handlers map[string]RequestHandler
	Routes   Routes
}

func NewMux() *Mux {
	return &Mux{
		Handlers: make(map[string]RequestHandler),
	}
}

func (m *Mux) Handle(r Request) {
	h, found := m.Handlers[r.RouteConfig]
	if !found {
		return
	}
	h(r)
}

func (m *Mux) Add(h RequestHandler, r RouteConfig) {
	m.Handlers[r.ID] = h
	m.Routes = append(m.Routes, r)
}
