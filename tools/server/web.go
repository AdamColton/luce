package server

import (
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type TemplateNames struct {
	SignIn       string
	Home         string
	HomeSignedIn string
}

func (s *Server) setRoutes(host string) {
	m := midware.New(
		midware.NewRedirect("Redirect"),
		s.Users.Midware(),
		midware.NewDecoder(valuedecoder.Form(), "Form"),
		midware.NewDecoder(valuedecoder.Query(), "URLData"),
	)

	r := s.coreserver.Router
	if host != "" {
		r = r.Host(host).Subrouter()
	}

	r.HandleFunc("/", m.Handle(s.home))
	r.HandleFunc("/user/signin", m.Handle(s.getSignIn)).Methods("GET")
	r.HandleFunc("/user/signin", m.Handle(s.postSignIn)).Methods("POST")
	r.HandleFunc("/user/signout", m.Handle(s.getSignOut)).Methods("GET")
	r.HandleFunc("/admin/users", m.Handle(s.adminUsers)).Methods("GET")
	r.HandleFunc("/services", m.Handle(s.listServices)).Methods("GET")
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
	n := s.TemplateNames.Home
	if u != nil {
		n = s.TemplateNames.HomeSignedIn
	}
	err := s.Templates.ExecuteTemplate(w, n, u)
	lerr.Panic(err)
}

var (
	servicesHeader = `<!DOCTYPE html>
<html>
	<head><title>Luce Server:Services</title></head>
	<body>`
	servicesFooter = `</body>
</html>`
)

func (s *Server) listServices(w http.ResponseWriter, r *http.Request, d *struct {
	Session *lusess.Session
}) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(servicesHeader)
	s.services.Each(func(name string, srv *serviceConn, done *bool) {
		sw.WriteStrings("<h2>", name, "</h2><ul>")
		for _, l := range srv.service.Links {
			sw.WriteStrings("<li><a href=\"", l.Get(s.coreserver.Host, s.coreserver.Addr), "\">", l.Name, "</a></li>")
		}

		sw.WriteString("</ul>")
	})
	sw.WriteString(servicesFooter)
}
