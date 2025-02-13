package server

import (
	"fmt"
	"html/template"

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
	core.Config
}

// Server runs a webserver.
type Server struct {
	coreserver *core.Server
	Users      *lusess.Store
	Templates  *template.Template
	lerr.ErrHandler
}

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
		Templates: c.Templates,
		ErrHandler: func(err error) {
			fmt.Println(err)
		},
	}

	return srv, nil
}

func (s *Server) Close() error {
	return s.coreserver.Close()
}

// Run the server. If Socket or ServiceSocket is defined, they will be run as
// well.
func (s *Server) Run() {
	s.coreserver.Run()
}
