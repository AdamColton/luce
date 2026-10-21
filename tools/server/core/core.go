package core

// == projects.Code.luce.server ==
// [ ] rename core
//	it only makes sense in context
//	should probably be core -> server
//  server -> lserver

import (
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/gorilla/mux"
)

// SSL holds the certificate and key paths used to serve over HTTPS. Leaving
// either empty serves plain HTTP.
type SSL struct {
	Cert, Key string
}

// Config holds the settings needed to build a Server.
type Config struct {
	Addr string
	// TODO: Host isn't used
	Host            string
	Socket          string
	CliStartMessage string
	SSL             SSL
}

// NewServer builds a Server from c, ready to Run.
func (c Config) NewServer() *Server {
	return &Server{
		Config:        c,
		Router:        mux.NewRouter(),
		httpserver:    &http.Server{},
		socketRunning: make(chan bool),
	}
}

// Server runs a gorilla/mux HTTP server, optionally alongside a unix-socket
// admin CLI driven by CliHandler.
type Server struct {
	Router *mux.Router
	Config
	CliHandler func(*cli.ExitClose) cli.Commander
	// OnShutdown is called by Shutdown, after the admin socket is closed and
	// before the HTTP server is.
	OnShutdown func()

	httpserver    *http.Server
	socket        *unixsocket.Socket
	socketRunning chan bool
}

// ListenAndServe starts the HTTP server, serving over TLS if Config.SSL has
// both a cert and a key. It blocks until the server stops.
func (s *Server) ListenAndServe() error {
	s.httpserver.Addr = s.Addr
	s.httpserver.Handler = s.Router
	if s.SSL.Cert != "" && s.SSL.Key != "" {
		return s.httpserver.ListenAndServeTLS(s.SSL.Cert, s.SSL.Key)
	}
	return s.httpserver.ListenAndServe()
}

// Close stops the HTTP server.
func (s *Server) Close() error {
	return s.httpserver.Close()
}

// Shutdown is a graceful Close: it stops the admin socket, which removes its
// socket file, calls OnShutdown if it is set, and then stops the HTTP server.
// Close alone leaves the socket file behind.
func (s *Server) Shutdown() error {
	if s.socket != nil {
		s.socket.Close()
	}
	if s.OnShutdown != nil {
		s.OnShutdown()
	}
	return s.Close()
}

// Run starts the admin socket, if configured, then blocks running the HTTP
// server. It panics on any ListenAndServe error other than
// http.ErrServerClosed.
func (s *Server) Run() {
	if s.Socket != "" && s.CliHandler != nil {
		go s.RunSocket()
	}

	lerr.Panic(s.ListenAndServe(), http.ErrServerClosed)
}
