package core

import (
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
)

func (s *Server) RunSocket() {
	unixsocket.CLISocket(s.Socket, s).Run()
}

func (s *Server) RunStdIO() {
	cli.StdIO(s)
}

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
