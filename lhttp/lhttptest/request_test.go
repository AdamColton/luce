package lhttptest_test

import (
	"net/url"
	"testing"

	"github.com/adamcolton/luce/lhttp/lhttptest"
	"github.com/stretchr/testify/assert"
)

type person struct {
	Name string
	Age  int
}

func TestRequestFromStruct(t *testing.T) {
	adam := person{Name: "Adam", Age: 37}
	req := lhttptest.NewRequest("/people", adam)

	get := req.GET()
	assert.Equal(t, "GET", get.Method)
	assert.Equal(t, "/people", get.URL.Path)
	assert.Equal(t, "Adam", get.URL.Query().Get("Name"))
	assert.Equal(t, "37", get.URL.Query().Get("Age"))

	post := req.POST()
	assert.Equal(t, "POST", post.Method)
	assert.Equal(t, "application/x-www-form-urlencoded", post.Header.Get("Content-Type"))
	assert.NoError(t, post.ParseForm())
	assert.Equal(t, "Adam", post.PostForm.Get("Name"))
	assert.Equal(t, "37", post.PostForm.Get("Age"))
}

func TestRequestFromValues(t *testing.T) {
	vals := url.Values{"q": {"go"}}
	req := lhttptest.NewRequest("/search", vals)
	assert.Equal(t, vals, req.Values)
	assert.Equal(t, "go", req.GET().URL.Query().Get("q"))
}

func TestRequestWithoutValues(t *testing.T) {
	req := lhttptest.NewRequest("/home", nil)
	assert.Nil(t, req.Values)

	get := req.GET()
	assert.Equal(t, "/home", get.URL.String())

	post := req.POST()
	assert.Empty(t, post.Header.Get("Content-Type"))
}

func TestEncodeValuesError(t *testing.T) {
	req := lhttptest.NewRequest("/", 42)
	assert.Empty(t, req.Values)
	assert.Error(t, req.EncodeValues(42))
}
