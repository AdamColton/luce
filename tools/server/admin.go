package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/adamcolton/luce/util/lusers/lusess"
)

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request, d *struct {
	Session  *lusess.Session
	Redirect string
}) {
	if !d.Session.User().In("admin") {
		d.Redirect = "/"
		return
	}

	fmt.Fprint(w, strings.Join(s.Users.List(), "<br>"))
}
