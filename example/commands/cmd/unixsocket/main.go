package main

import (
	"net"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/example/commands/logic"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
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
		rw := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReadWriter(conn)

		rwPipe, _, _ := channel.NewPipe(rw.In, rw.Out)
		prefixPipe := packeter.Run(prefix.New[uint32](), rwPipe)

		w := channel.Writer{prefixPipe.Out}

		r.Context = cli.NewContext(w, prefixPipe.In, nil)
		r.Run()
	})

	s.Run()
}
