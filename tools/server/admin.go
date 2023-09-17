package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/adamcolton/luce/util/lusers/lusess"
)

func (s *Server) adminUsers(w http.ResponseWriter, r *http.Request, d *struct {
	Session *lusess.Session
}) {
	if !d.Session.User().In("admin") {
		redirect(w, r)
	}

	fmt.Fprint(w, strings.Join(s.Users.List(), "<br>"))
}
