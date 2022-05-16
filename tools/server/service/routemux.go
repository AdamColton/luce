package service

// RequestHandler will handle a request but it will not respond to it.
type RequestHandler func(r Request)

// RequestResponder will hand a request and return a reponse.
type RequestResponder func(r Request) Response

// Mux maps the Requests to their handlers.
type Mux struct {
	Handlers map[string]RequestHandler
	Routes   Routes
}

// NewMux creates a Mux.
func NewMux() *Mux {
	return &Mux{
		Handlers: make(map[string]RequestHandler),
	}
}

// Handle a request. If there is no handler, nothing happens.
func (m *Mux) Handle(r Request) {
	h, found := m.Handlers[r.RouteConfig]
	if !found {
		return
	}
	h(r)
}

// Add a RequestHandler to the Mux mapped to the RouteConfig.
func (m *Mux) Add(h RequestHandler, r RouteConfig) {
	m.Handlers[r.ID] = h
	m.Routes = append(m.Routes, r)
}
