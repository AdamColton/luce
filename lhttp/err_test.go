package lhttp_test

import (
	"net/http"
	"testing"

	"github.com/adamcolton/luce/lhttp"
	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	var err error = lhttp.Err{Code: http.StatusTeapot, String: "short and stout"}
	assert.EqualError(t, err, "short and stout")
	assert.Equal(t, http.StatusTeapot, lhttp.ErrStatus(err))
}
