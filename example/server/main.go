package main

import (
	"encoding/gob"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/quasoft/memstore"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lhttp/tokens"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

func main() {
	fac := ephemeral.Factory(bytebtree.New, 10)
	s := &Server{
		Router: mux.NewRouter(),
		Addr:   ":6060",
		Users: &lusess.Store{
			Store: memstore.NewMemStore(
				[]byte("authkey123"),
				[]byte("enckey12341234567890123456789012"),
			),
			UserStore: lusers.MustUserStore(fac),
			FieldName: "Session",
		},
		tokens: tokens.New(time.Second * 5),
		server: &http.Server{},
	}
	s.loadTemplates()
	s.routes()

	s.initUsers()

	gob.Register((*lusers.User)(nil))

	s.server.Handler = s.Router
	s.server.Addr = s.Addr
	go s.server.ListenAndServe()
	s.RunSocket()
}

func (s *Server) initUsers() {
	us := s.Users.UserStore

	au, _ := us.Create("adminUser", "test")
	a, _ := us.Group("admin")
	a.AddUser(au)
	us.Update(au)

	eu, _ := us.Create("editorUser", "test")
	e, _ := us.Group("editor")
	e.AddUser(eu)
	us.Update((eu))

	us.Group("subscriber")
}

func (s *Server) loadTemplates() {
	s.Tmpls = template.Must(template.ParseGlob("*.html"))
}

type Server struct {
	Router   *mux.Router
	Addr     string
	Users    *lusess.Store
	Tmpls    *template.Template
	Settings struct {
		AdminLockUserCreation bool
	}
	tokens *tokens.Tokens
	server *http.Server
}
