package service

import (
	"net"

	"github.com/adamcolton/luce/lerr"
)

// Client is a service's connection to the server: a Mux to route incoming
// requests to handlers, over a Conn.
type Client struct {
	*Mux
	*Conn
}

// NewClient connects to the server's service socket at addr.
func NewClient(addr string) (*Client, error) {
	netConn, err := net.Dial("unix", addr)
	if err != nil {
		return nil, err
	}

	conn, err := NewConn(netConn)
	if err != nil {
		return nil, err
	}

	mux := NewMux()
	err = conn.Listener.RegisterHandlers(mux.Handle)
	if err != nil {
		return nil, err
	}

	return &Client{
		Mux:  mux,
		Conn: conn,
	}, nil
}

// Run sends the Client's Service registration, then blocks handling
// incoming requests.
func (c *Client) Run() {
	c.Sender.Send(c.Mux.Service)
	c.Listener.Run()
}

// Add registers route, handled by h, with the Client's Service and Mux.
func (c *Client) Add(h RequestResponder, route *Route) {
	lerr.Panic(route.Validate())
	c.Service.Routes = append(c.Service.Routes, *route)
	fn := func(r *Request) {
		err := c.Sender.Send(h(r))
		// == projects.Code.luce.server ==
		// [ ] Add's response-send failure panics the handler goroutine
		//  a failed Sender.Send panics rather than reporting the error,
		//  which takes down the whole client on one bad send
		lerr.Panic(err)
	}
	c.Mux.Handlers[route.ID] = fn
}

// Handle is a chainable helper. It calls c.Add with r and h.
func (r *Route) Handle(c *Client, h RequestResponder) {
	c.Add(h, r)
}
