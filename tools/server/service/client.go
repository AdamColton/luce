package service

import (
	"net"

	"github.com/adamcolton/luce/lerr"
)

// Client is a sub-service running on a luce Server.
type Client struct {
	*Mux
	*Conn
}

// MustClient creates a client and will panic if there are any errors. The addr
// value is the unix socket for communicating with the server.
func MustClient(addr string) *Client {
	c, err := NewClient(addr)
	lerr.Panic(err)
	return c
}

// NewClient creates a client. The addr value is the unix socket for
// communicating with the server.
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
		conn.NetConn.Close()
		return nil, err
	}

	return &Client{
		Mux:  mux,
		Conn: conn,
	}, nil
}

// Run sends the routes to the server and runs the listener.
func (c *Client) Run() {
	c.Sender.Send(c.Mux.Routes)
	c.Listener.Run()
}

// Add a Route, RequestResponder is the handler and RouteConfig describes the
// route.
func (c *Client) Add(h RequestResponder, r RouteConfig) {
	lerr.Panic(r.Validate())
	fn := func(r Request) {
		c.Sender.Send(h(r))
	}
	c.Mux.Add(fn, r)
}
