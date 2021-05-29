package server

import (
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type signinData struct {
	Message string
}

func (signinData) Title() string { return "Sign In" }

func (s *Server) getSignIn(w http.ResponseWriter, r *http.Request, d *struct {
	Session  *lusess.Session
	Redirect string
	URLData  *signinData
}) {
	if d.Session.User() != nil {
		d.Redirect = "/"
		return
	}
	err := s.Templates.ExecuteTemplate(w, s.TemplateNames.SignIn, d.URLData)
	lerr.Panic(err)
}

func (s *Server) postSignIn(w http.ResponseWriter, r *http.Request, d *struct {
	Form     *lusess.Login
	Session  *lusess.Session
	Redirect string
}) {
	_, err := d.Session.Login(d.Form)
	if err != nil {
		d.Redirect = "/user/signin?Message=login+failed"
		return
	}
	d.Redirect = "/"
}

func (s *Server) getSignOut(w http.ResponseWriter, r *http.Request, d *struct {
	Session  *lusess.Session
	Redirect string
}) {
	d.Session.Logout()
	d.Redirect = "/"
}
