package server

import (
	"fmt"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/formdecoder"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type TemplateNames struct {
	SignIn       string
	Home         string
	HomeSignedIn string
}

func (s *Server) setRoutes(host string) {
	m := midware.NewMagic(
		s.Users,
		midware.NewDecoder(formdecoder.New(), "Form"),
		midware.Redirect{"Redirect"},
	)
	r := s.Router
	if host != "" {
		r = r.Host(host).Subrouter()
	}
	r.HandleFunc("/setCookie", setCookie)
	r.HandleFunc("/", m.Handle(s.home))
	r.HandleFunc("/user/signin", m.Handle(s.getSignIn)).Methods("GET")
	r.HandleFunc("/user/signin", m.Handle(s.postSignIn)).Methods("POST")
	r.HandleFunc("/user/signout", m.Handle(s.getSignOut)).Methods("GET")
	r.HandleFunc("/admin/users", m.Handle(s.adminUsers)).Methods("GET")
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "testing",
		Value: "this is a test",
	}
	http.SetCookie(w, cookie)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request, d *struct {
	Session *lusess.Session
}) {
	u := d.Session.User()
	fmt.Println(u)
	n := s.TemplateNames.Home
	if u != nil {
		n = s.TemplateNames.HomeSignedIn
	}
	err := s.Templates.ExecuteTemplate(w, n, u)
	lerr.Panic(err)
}
