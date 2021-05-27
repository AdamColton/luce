package tokens_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adamcolton/luce/lhttp/tokens"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAndCall(t *testing.T) {
	ts := tokens.New(time.Minute)

	var called bool
	token, toqToken := ts.Register(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	assert.NotEmpty(t, token)
	assert.NotNil(t, toqToken)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ts.Call(token, w, r)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterNilFunc(t *testing.T) {
	ts := tokens.New(time.Minute)
	token, toqToken := ts.Register(nil)
	assert.Empty(t, token)
	assert.Nil(t, toqToken)
}

func TestRegisterExpires(t *testing.T) {
	ts := tokens.New(5 * time.Millisecond)
	calls := 0
	token, _ := ts.Register(func(w http.ResponseWriter, r *http.Request) {
		calls++
	})

	call := func() {
		ts.Call(token, httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	}

	call()
	assert.Equal(t, 1, calls)

	// well past the timeout: the toq callback has removed the token
	time.Sleep(50 * time.Millisecond)
	call()
	assert.Equal(t, 1, calls)
}

func TestCallUnknownToken(t *testing.T) {
	ts := tokens.New(time.Minute)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	// an unregistered token is a no-op, not a panic
	ts.Call("no-such-token", w, r)
	assert.Equal(t, 200, w.Code)
}

func TestPost(t *testing.T) {
	ts := tokens.New(time.Minute)
	var got string
	token, _ := ts.Register(func(w http.ResponseWriter, r *http.Request) {
		got = "called"
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(token))
	ts.Post(w, r)

	assert.Equal(t, "called", got)
}

func TestGet(t *testing.T) {
	ts := tokens.New(time.Minute)
	var got string
	token, _ := ts.Register(func(w http.ResponseWriter, r *http.Request) {
		got = "called"
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/callback/"+token, nil)
	ts.Get(w, r)

	assert.Equal(t, "called", got)
}
