package server

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

// Config is used to create a Server. It holds all the configuration values.
// This follows the builder pattern instead of using a function that would
// take many arguments.
type Config struct {
	SessionStore sessions.Store
	Templates    *template.Template
	UserStore    store.Factory
	TemplateNames
	Addr          string
	Socket        string
	ServiceSocket string
	Host          string
}

// Server runs a webserver.
type Server struct {
	Router        *mux.Router
	Addr          string
	Users         *lusess.Store
	Settings      Settings
	Templates     *template.Template
	Socket        string
	ServiceSocket string
	TemplateNames
	server        *http.Server
	serviceRoutes map[string]*serviceRoute
	lerr.ErrHandler
}

var TimeoutDuration = time.Second * 5

// New Server using the values from the Config.
func (c *Config) New() (*Server, error) {
	us, err := lusers.NewUserStore(c.UserStore)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		Router: mux.NewRouter(),
		Users: &lusess.Store{
			Store:     c.SessionStore,
			UserStore: us,
			Decoder:   schema.NewDecoder(),
			FieldName: "Session",
		},
		Addr:          c.Addr,
		Templates:     c.Templates,
		Socket:        c.Socket,
		ServiceSocket: c.ServiceSocket,
		TemplateNames: c.TemplateNames,
		server:        &http.Server{},
		serviceRoutes: make(map[string]*serviceRoute),
		ErrHandler: func(err error) {
			fmt.Println(err)
		},
	}

	srv.setRoutes(c.Host)
	return srv, nil
}

func (s *Server) ListenAndServe() error {
	s.server.Addr = s.Addr
	s.server.Handler = s.Router
	return s.server.ListenAndServe()
}

func (s *Server) Close() error {
	return s.server.Close()
}

// Run the server. If Socket or ServiceSocket is defined, they will be run as
// well.
func (s *Server) Run() {
	if s.Socket != "" {
		go func() {
			s.RunSocket()
		}()
	}

	if s.ServiceSocket != "" {
		go func() {
			s.RunServiceSocket()
		}()
	}

	lerr.Panic(s.ListenAndServe(), http.ErrServerClosed)
}
