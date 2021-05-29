package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/adamcolton/luce/ds/toq"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

type Server struct {
	Router    *mux.Router
	Addr      string
	Users     *lusess.Store
	Templates *template.Template
	Settings  struct {
		AdminLockUserCreation bool
	}
	tokens map[string]http.HandlerFunc
	toq    *toq.TimeoutQueue
	UserTemplates
	server *http.Server
}

var TimeoutDuration = time.Second * 5

func New(ses sessions.Store, fac store.Factory, t *template.Template, ut UserTemplates) (*Server, error) {
	us, err := lusers.NewUserStore(fac)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		Router: mux.NewRouter(),
		Users: &lusess.Store{
			Store:     ses,
			UserStore: us,
			Decoder:   schema.NewDecoder(),
		},
		Templates:     t,
		tokens:        make(map[string]http.HandlerFunc),
		toq:           toq.New(TimeoutDuration, 10),
		UserTemplates: ut,
		server:        &http.Server{},
	}

	srv.setRoutes()
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
