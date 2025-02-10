package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := []byte("test body")
	resp := req.Response(body)
	assert.Equal(t, body, resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, service.ResponseTypeID32, resp.TypeID32())
}

func TestResponseString(t *testing.T) {
	req := &service.Request{
		ID: 31415,
	}
	body := "test body"
	resp := req.ResponseString(body)
	assert.Equal(t, []byte(body), resp.Body)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, service.ResponseTypeID32, resp.TypeID32())
}
