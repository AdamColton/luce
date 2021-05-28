package lusess

import (
	"encoding/gob"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/gorilla/sessions"
)

var (
	// StoreName is the sessions.Store name under which the user's session
	// is kept.
	StoreName = "User"
	// ValueName is the key, within that session, that holds the logged-in
	// *lusers.User.
	ValueName = "User"
)

const (
	// ErrLoginFailed is returned when a username/password pair does not
	// match a stored User.
	ErrLoginFailed = lerr.Str("Login failed")
)

func init() {
	gob.Register((*lusers.User)(nil))
	gob.Register((*lusers.Group)(nil))
}

// Store combines a gorilla sessions.Store with a lusers.UserStore and a form
// Decoder, so it can both authenticate against stored Users and manage
// sessions for them. Decoder is used by Login to parse POSTed form values
// into a Login. FieldName names the field Store.InitilizeField fills when
// Store is used as a midware.
type Store struct {
	sessions.Store
	*lusers.UserStore
	Decoder interface {
		Decode(interface{}, map[string][]string) error
	}
	FieldName string
}

// Session wraps a gorilla sessions.Session for one request, keeping the
// ResponseWriter and Request needed to Save it.
type Session struct {
	*sessions.Session
	Store *Store
	W     http.ResponseWriter
	R     *http.Request
}

// Session fetches r's session named StoreName and wraps it, ready to read
// or set the logged-in User and to be Saved.
func (s *Store) Session(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sess, err := s.Get(r, StoreName)
	if err != nil {
		return nil, err
	}
	return &Session{
		Session: sess,
		Store:   s,
		W:       w,
		R:       r,
	}, nil
}

// Login decodes r's POST form into a Login, authenticates it against s's
// UserStore, and, on success, sets the logged-in User on r's Session and
// saves it. On failure the session is left unchanged and not saved.
func (s *Store) Login(w http.ResponseWriter, r *http.Request) (*Session, error) {
	err := r.ParseForm()
	lerr.Panic(err)

	var login Login
	err = s.Decoder.Decode(&login, r.PostForm)
	if err != nil {
		return nil, err
	}

	sess, err := s.Session(w, r)
	if err != nil {
		return nil, err
	}

	_, err = sess.Login(&login)
	if err != nil {
		return sess, err
	}
	return sess, sess.Save()
}

// Login holds the form fields Store.Login decodes a request into.
type Login struct {
	Username, Password string
}

// Login authenticates l against s.Store's UserStore and, on success, sets
// the resulting User on the session. It returns ErrLoginFailed if the user
// does not exist or the password is wrong.
func (s *Session) Login(l *Login) (*lusers.User, error) {
	u, err := s.Store.UserStore.Login(l.Username, l.Password)
	if err != nil || u == nil {
		return nil, ErrLoginFailed
	}

	s.Session.Values[ValueName] = u
	return u, nil
}

// Save persists the session, writing its cookie to W.
func (s *Session) Save() error {
	return s.Session.Save(s.R, s.W)
}

// User returns the session's logged-in User, or nil if none is set.
func (s *Session) User() *lusers.User {
	i := s.Session.Values[ValueName]
	if i == nil {
		return nil
	}
	u, _ := i.(*lusers.User)
	return u
}

// SetUser sets the session's logged-in User directly, without going
// through Login.
func (s *Session) SetUser(u *lusers.User) {
	s.Session.Values[ValueName] = u
}

// Logout clears the session's logged-in User.
func (s *Session) Logout() {
	delete(s.Session.Values, ValueName)
}
