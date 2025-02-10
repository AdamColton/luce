package service_test

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestRequestServeFile(t *testing.T) {
	old := service.OS
	defer func() { service.OS = old }()
	service.OS = lfilemock.Parse(map[string]any{
		"hello.txt": "hello world",
	}).Repository()

	req := &service.Request{ID: 31415}
	resp := req.ServeFile("hello.txt")
	assert.Equal(t, []byte("hello world"), resp.Body)

	resp = req.ServeFile("missing.txt")
	assert.NotEqual(t, []byte("hello world"), resp.Body)
}
