package core

import (
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/gorilla/mux"
)

type SSL struct {
	Cert, Key string
}

type Config struct {
	Addr            string
	Host            string
	Socket          string
	CliStartMessage string
	SSL             SSL
}

func (c Config) NewServer() *Server {
	return &Server{
		Config:        c,
		Router:        mux.NewRouter(),
		httpserver:    &http.Server{},
		socketRunning: make(chan bool),
	}
}

type Server struct {
	Router *mux.Router
	Config
	CliHandler func(*cli.ExitClose) cli.Commander

	httpserver    *http.Server
	socket        *unixsocket.Socket
	socketRunning chan bool
}

func (s *Server) ListenAndServe() error {
	s.httpserver.Addr = s.Addr
	s.httpserver.Handler = s.Router
	if s.SSL.Cert != "" && s.SSL.Key != "" {
		return s.httpserver.ListenAndServeTLS(s.SSL.Cert, s.SSL.Key)
	}
	return s.httpserver.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpserver.Close()
}

func (s *Server) Run() {
	if s.Socket != "" && s.CliHandler != nil {
		go s.RunSocket()
	}

	lerr.Panic(s.ListenAndServe(), http.ErrServerClosed)
}
