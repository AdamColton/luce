package main

import (
	"net"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/unixsocket"
)

func main() {
	var s *unixsocket.Socket
	s = unixsocket.New("/tmp/testsocket.sock", func(conn net.Conn) {
		ec := cli.NewExitClose(
			func() { conn.Close() },
			func() { s.Close() },
		)
		ho := logic.New(ec)
		r := ho.Runner()
		rdr := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReader(conn)

		r.Context = cli.NewContext(conn, rdr.Out, nil)
		r.Run()
	})

	s.Run()
}
