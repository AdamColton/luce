package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestServiceRoutePointer(t *testing.T) {
	// s := &service.Service{}
	// r := s.New("foo")
	// r.Method = "GET"
	// assert.Equal(t, "GET", s.Routes[0].Method)

	// service.NewServiceRoute("test").
	// 	WithBody().
	// 	RequireGroup(mediaAdmin).
	// 	Post().
	// 	Handle(conn, s.Cmd)

	c := &service.Client{}
	c.Mux = &service.Mux{}
	c.Service = &service.Service{}
	c.Mux.Handlers = make(map[string]service.RequestHandler)
	r := service.NewRoute("foo")
	assert.Equal(t, "/foo", r.Path)
	h := func(r *service.Request) *service.Response {
		return nil
	}
	r.Handle(c, h)

}
