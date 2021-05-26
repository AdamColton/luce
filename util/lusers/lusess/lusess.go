package lusess

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/gorilla/mux"
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
	Router  *mux.Router
	Decoder interface {
		Decode(interface{}, map[string][]string) error
	}
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

type HandlerFunc func(http.ResponseWriter, *http.Request, *Session)

func (s *Store) HandlerFunc(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.Session(w, r)
		if err != nil {
			log.Println(err)
			return
		}
		fn(w, r, sess)
	}
}

func (s *Store) HandleSession(path string, fn HandlerFunc) *mux.Route {
	return s.Router.HandleFunc(path, s.HandlerFunc(fn))
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
	u, err := s.Store.GetByName(l.Username)
	if err != nil || u == nil || u.CheckPassword(l.Password) != nil {
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
