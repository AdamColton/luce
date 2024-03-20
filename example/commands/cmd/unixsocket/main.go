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
		exit := func() { conn.Close() }
		cls := func() { s.Close() }
		r := logic.RWHandler(exit, cls)
		in, _ := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReader(conn)

		r.Context = cli.NewContext(conn, in, nil)
		r.Run()
	})

	s.Run()
}
