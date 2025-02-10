package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestResponseSetHeader(t *testing.T) {
	resp := &service.Response{}
	resp.SetHeader("X-Test", "1")
	assert.Equal(t, "1", resp.Header.Get("X-Test"))
}

func TestResponseContentType(t *testing.T) {
	resp := &service.Response{}
	resp.ContentType("text/plain")
	assert.Equal(t, "text/plain", resp.Header.Get(service.ContentType))
}
