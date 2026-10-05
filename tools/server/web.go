package server

import (
	"cmp"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

// TemplateNames names the templates Server renders for the built-in user
// pages.
type TemplateNames struct {
	SignIn       string
	Home         string
	HomeSignedIn string
	Services     string
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
	r.HandleFunc("/admin/listBashCmds", m.Handle(s.listBashCmds)).Methods("GET")

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

var ltStr = cmp.Less[string]

// serviceLink is a single link rendered on the services page, with its URL
// already resolved against the server's host and port.
type serviceLink struct {
	Name, URL string
}

// serviceListing is one service's name and resolved links, as rendered on
// the services page.
type serviceListing struct {
	Name  string
	Links []serviceLink
}

func (s *Server) listServices(w http.ResponseWriter, r *http.Request, d *struct {
	Session *lusess.Session
}) {
	var listings []serviceListing
	s.services.SortedEachKey(ltStr, nil, func(name string, srv *serviceConn, done *bool) {
		links := make([]serviceLink, len(srv.service.Links))
		for i, l := range srv.service.Links {
			links[i] = serviceLink{Name: l.Name, URL: l.Get(s.coreserver.Host, s.coreserver.Addr)}
		}
		listings = append(listings, serviceListing{Name: name, Links: links})
	})
	err := s.Templates.ExecuteTemplate(w, s.TemplateNames.Services, listings)
	lerr.Panic(err)
}
