package core

import (
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
)

// RunSocket starts the admin unix socket and blocks until it stops.
func (s *Server) RunSocket() {
	unixsocket.CLISocket(s.Socket, s).Run()
}

// RunStdIO drives the admin CLI over the process's stdin/stdout.
func (s *Server) RunStdIO() {
	cli.StdIO(s)
}

// Cli builds the CliHandler's command tree for one connection and runs it
// against ctx, wiring onExit and the HTTP server's Close as the
// connection's exit and close actions.
func (s *Server) Cli(ctx cli.Context, onExit func()) {
	onClose := func() {
		s.Close()
	}
	ec := cli.NewExitClose(onExit, onClose)
	c := s.CliHandler(ec)

	r := cli.NewRunner(c, ctx)
	r.StartMessage = s.CliStartMessage
	r.Run()
}
