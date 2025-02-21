package service

import (
	"net"

	"github.com/adamcolton/luce/lerr"
)

type Client struct {
	*Mux
	*Conn
}

// Name   string
// 	Host   string
// 	Base   string

func NewService(name, host, base, addr string) (*Client, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	c.Service.Name = name
	c.Service.Base = base
	if host != "" {
		c.Service.Host = host + ".{domain:.*}"
	}
	return c, nil
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
	c.Sender.Send(c.Mux.Service)
	c.Listener.Run()
}

func (c *Client) Add(h RequestResponder, route *Route) {
	lerr.Panic(route.Validate())
	c.Service.Routes = append(c.Service.Routes, *route)
	fn := func(r *Request) {
		err := c.Sender.Send(h(r))
		// TODO: handle error
		lerr.Panic(err)
	}
	c.Mux.Handlers[route.ID] = fn
}

func (r *Route) Handle(c *Client, h RequestResponder) {
	c.Add(h, r)
}
