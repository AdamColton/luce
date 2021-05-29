package server

import (
	"fmt"
	"net/http"

	"github.com/adamcolton/luce/lerr"
)

type UserTemplates struct {
	SignIn string
}

type signinData struct {
	Message string
}

func (signinData) Title() string { return "Sign In" }

func (s *Server) getSignIn(w http.ResponseWriter, r *http.Request) {
	err := s.Templates.ExecuteTemplate(w, s.UserTemplates.SignIn, signinData{})
	lerr.Panic(err)
}

func (s *Server) postSignIn(w http.ResponseWriter, r *http.Request) {
	sess, err := s.Users.Login(w, r)
	lerr.Panic(err)

	fmt.Fprint(w, sess.User())
}
