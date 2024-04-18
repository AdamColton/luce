package lhttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adamcolton/luce/lhttp"
	"github.com/stretchr/testify/assert"
)

func TestErrHandlerCheck(t *testing.T) {
	var handled error
	h := lhttp.ErrHandler(func(w http.ResponseWriter, r *http.Request, err error) {
		handled = err
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	assert.False(t, h.Check(w, r, nil))
	assert.NoError(t, handled)

	boom := errors.New("boom")
	assert.True(t, h.Check(w, r, boom))
	assert.Equal(t, boom, handled)

	var none lhttp.ErrHandler
	assert.True(t, none.Check(w, r, boom))
	assert.False(t, none.Check(w, r, nil))
}

type statusError struct{ status int }

func (e statusError) Error() string { return "status error" }
func (e statusError) Status() int   { return e.status }

func TestErrStatus(t *testing.T) {
	assert.Equal(t, 0, lhttp.ErrStatus(nil))
	assert.Equal(t, http.StatusInternalServerError, lhttp.ErrStatus(errors.New("plain")))
	assert.Equal(t, http.StatusNotFound, lhttp.ErrStatus(statusError{http.StatusNotFound}))
}
