package server

import (
	"net/http"
)

func (s *Server) setRoutes() {
	r := s.Router

	r.HandleFunc("/", s.sayHi)
	r.HandleFunc("/user/signin", s.getSignIn).Methods("GET")
	r.HandleFunc("/user/signin", s.postSignIn).Methods("POST")
}

func (s *Server) sayHi(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hi!"))
}
