package service

import (
	"net"

	"github.com/adamcolton/luce/lerr"
)

type Client struct {
	*Mux
	*Conn
}

func MustClient(addr string) *Client {
	c, err := NewClient(addr)
	lerr.Panic(err)
	return c
}

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

func (c *Client) Run() {
	c.Sender.Send(c.Mux.Routes)
	c.Listener.Run()
}

func (c *Client) Add(h RequestResponder, r RouteConfig) {
	lerr.Panic(r.Validate())
	fn := func(r Request) {
		c.Sender.Send(h(r))
	}
	c.Mux.Add(fn, r)
}
