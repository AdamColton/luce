package lusess

import (
	"encoding/gob"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/gorilla/sessions"
)

var (
	StoreName = "User"
	ValueName = "User"
)

const (
	ErrLoginFailed = lerr.Str("Login failed")
)

func init() {
	gob.Register((*lusers.User)(nil))
	gob.Register((*lusers.Group)(nil))
}

type Store struct {
	sessions.Store
	*lusers.UserStore
	Decoder interface {
		Decode(interface{}, map[string][]string) error
	}
	FieldName string
}

type Session struct {
	*sessions.Session
	Store *Store
	W     http.ResponseWriter
	R     *http.Request
}

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

func (s *Store) User(r *http.Request) (*lusers.User, error) {
	sess, err := s.Get(r, StoreName)
	if err != nil {
		return nil, err
	}

	i := sess.Values[ValueName]
	if i == nil {
		return nil, nil
	}
	u, _ := i.(*lusers.User)
	return u, nil

}

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
		err = sess.Save()
	}
	return sess, err
}

type Login struct {
	Username, Password string
}

func (s *Session) Login(l *Login) (*lusers.User, error) {
	u, err := s.Store.UserStore.Login(l.Username, l.Password)
	if err != nil || u == nil {
		return nil, ErrLoginFailed
	}

	s.Session.Values[ValueName] = u
	return u, nil
}

func (s *Session) Save() error {
	return s.Session.Save(s.R, s.W)
}

func (s *Session) User() *lusers.User {
	i := s.Session.Values[ValueName]
	if i == nil {
		return nil
	}
	u, _ := i.(*lusers.User)
	return u
}

func (s *Session) SetUser(u *lusers.User) {
	s.Session.Values[ValueName] = u
}

func (s *Session) Logout() {
	delete(s.Session.Values, ValueName)
}
