package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestLinkGet(t *testing.T) {
	local := service.Link{Name: "home", Path: "/home"}
	assert.Equal(t, "/home", local.Get("example.com", ":8080"))

	remote := service.Link{Name: "api", Host: "api", Path: "/v1"}
	assert.Equal(t, "https://api.example.com:8080/v1", remote.Get("example.com", ":8080"))
}

func TestServiceValidate(t *testing.T) {
	s := &service.Service{
		Routes: []service.Route{
			*service.NewRoute("foo"),
			{},
		},
	}
	err := s.Validate()
	assert.Error(t, err)

	s = &service.Service{Routes: []service.Route{*service.NewRoute("foo")}}
	assert.NoError(t, s.Validate())
}

func TestServiceAddLink(t *testing.T) {
	s := &service.Service{Base: "/widget"}
	s.AddLink("list", "", "list")
	assert.Len(t, s.Links, 1)
	assert.Equal(t, "list", s.Links[0].Name)
	assert.Equal(t, "/widget/list", s.Links[0].Path)
}
