package server

import (
	"fmt"
	"html/template"
	"time"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/tools/server/core"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

// Config is used to create a Server. It holds all the configuration values.
// This follows the builder pattern instead of using a function that would
// take many arguments.
type Config struct {
	SessionStore sessions.Store
	Templates    *template.Template
	UserStore    store.NestedFactory
	TemplateNames
	ServiceSocket string
	core.Config
}

// Server runs a webserver.
type Server struct {
	coreserver    *core.Server
	Users         *lusess.Store
	Settings      Settings
	Templates     *template.Template
	ServiceSocket string
	TemplateNames
	// TODO: make this safe map
	serviceRoutes lmap.Wrapper[string, *serviceRoute]
	services      lmap.Wrapper[string, *serviceConn]
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
		coreserver: c.Config.NewServer(),
		Users: &lusess.Store{
			Store:     c.SessionStore,
			UserStore: us,
			Decoder:   schema.NewDecoder(),
			FieldName: "Session",
		},
		Templates:     c.Templates,
		ServiceSocket: c.ServiceSocket,
		TemplateNames: c.TemplateNames,
		serviceRoutes: lmap.NewSafe[string, *serviceRoute](nil),
		services:      lmap.NewSafe[string, *serviceConn](nil),
		ErrHandler: func(err error) {
			fmt.Println(err)
		},
	}
	srv.coreserver.CliHandler = srv.coreCommander

	srv.setRoutes(c.Host)
	return srv, nil
}

func (s *Server) Close() error {
	return s.coreserver.Close()
}

// Run the server. If Socket or ServiceSocket is defined, they will be run as
// well.
func (s *Server) Run() {
	if s.ServiceSocket != "" {
		go func() {
			lerr.Panic(s.RunServiceSocket())
		}()
	}

	s.coreserver.Run()
}
