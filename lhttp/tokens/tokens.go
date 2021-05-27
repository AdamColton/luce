package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/adamcolton/luce/ds/toq"
)

// Tokens can be used to create temporary callbacks.
type Tokens struct {
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
func (t Tokens) Register(fn http.HandlerFunc) (string, toq.Token) {
	if fn == nil {
		return "", nil
	}
	b := make([]byte, 10)
	rand.Read(b)
	token := base64.RawURLEncoding.EncodeToString(b)
	t.tokens[token] = fn
	toqToken := t.toq.Add(func() {
		delete(t.tokens, token)
	})
	return token, toqToken
}

// Call the http.HandlerFunc associated with the provided token.
func (t Tokens) Call(token string, w http.ResponseWriter, r *http.Request) {
	fn := t.tokens[token]
	if fn != nil {
		fn(w, r)
	}
}

// Post reads the body of the request as a token and invoke Tokens.Call.
func (t Tokens) Post(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	t.Call(string(b), w, r)
}

// Get assumes the last portion of a url is the token
func (t Tokens) Get(w http.ResponseWriter, r *http.Request) {
	token := path.Base(r.URL.Path)
	t.Call(token, w, r)
}
