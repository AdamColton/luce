package main

import (
	"net"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
)

func main() {
	var s *unixsocket.Socket
	s = unixsocket.New("/tmp/testsocket.sock", func(conn net.Conn) {
		pipe := unixsocket.ConnPipe(conn)
		w := channel.Writer{pipe.Snd}

		ec := cli.NewExitClose(
			func() { conn.Close() },
			func() { s.Close() },
		)
		ho := logic.New(ec)
		r := ho.Runner()
		r.Context = cli.NewContext(w, pipe.Rcv, nil)
		r.Run()
	})

	s.Run()
}
