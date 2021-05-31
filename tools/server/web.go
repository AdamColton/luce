package server

import (
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

func (s *Server) setRoutes() {
	m := midware.NewMagic(
		s.Users,
		midware.NewDecoder(formdecoder.New(), "Form"),
	)
	r := s.Router

	r.HandleFunc("/", m.Handle(s.home))
	r.HandleFunc("/user/signin", m.Handle(s.getSignIn)).Methods("GET")
	r.HandleFunc("/user/signin", m.Handle(s.postSignIn)).Methods("POST")
	r.HandleFunc("/user/signout", m.Handle(s.getSignOut)).Methods("GET")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request, d *struct {
	Session *lusess.Session
}) {
	u := d.Session.User()
	n := s.TemplateNames.Home
	if u != nil {
		n = s.TemplateNames.HomeSignedIn
	}
	err := s.Templates.ExecuteTemplate(w, n, u)
	lerr.Panic(err)
}
