// Package tokens hands out short-lived, random tokens that map to an
// http.HandlerFunc, so a URL or form field can carry a one-time callback
// instead of the handler itself. Each token expires and is forgotten after
// its Tokens' timeout.
package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/adamcolton/luce/ds/toq"
)

// Tokens can be used to create temporary callbacks.
type Tokens struct {
	mu     sync.Mutex
	tokens map[string]http.HandlerFunc
	toq    *toq.TimeoutQueue
}

// New initilizes an instance of Tokens with given duration.
func New(timeout time.Duration) *Tokens {
	return &Tokens{
		tokens: make(map[string]http.HandlerFunc),
		toq:    toq.New(timeout, 20),
	}
}

// Register an http.HandlerFunc. The string is the token. The toq.Token provides
// control over the timeout operation.
func (t *Tokens) Register(fn http.HandlerFunc) (string, toq.Token) {
	if fn == nil {
		return "", nil
	}
	b := make([]byte, 10)
	rand.Read(b)
	token := base64.RawURLEncoding.EncodeToString(b)
	t.mu.Lock()
	t.tokens[token] = fn
	t.mu.Unlock()
	toqToken := t.toq.Add(func() {
		t.mu.Lock()
		delete(t.tokens, token)
		t.mu.Unlock()
	})
	return token, toqToken
}

// Call the http.HandlerFunc associated with the provided token.
func (t *Tokens) Call(token string, w http.ResponseWriter, r *http.Request) {
	t.mu.Lock()
	fn := t.tokens[token]
	t.mu.Unlock()
	if fn != nil {
		fn(w, r)
	}
}

// Post reads the body of the request as a token and invoke Tokens.Call.
func (t *Tokens) Post(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	t.Call(string(b), w, r)
}

// Get assumes the last portion of a url is the token
func (t *Tokens) Get(w http.ResponseWriter, r *http.Request) {
	token := path.Base(r.URL.Path)
	t.Call(token, w, r)
}
