package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestMuxHandle(t *testing.T) {
	m := service.NewMux()
	called := false
	m.Handlers["found"] = func(r *service.Request) {
		called = true
	}

	m.Handle(&service.Request{RouteConfig: "not-found"})
	assert.False(t, called)

	m.Handle(&service.Request{RouteConfig: "found"})
	assert.True(t, called)
}
