package unixsocket

import (
	"io"
	"net"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
	"github.com/adamcolton/luce/util/reflector"
)

func ConnPipe(conn io.ReadWriter) channel.Pipe[[]byte] {
	rw := iobus.Config{
		CloseOnEOF: true,
		Sleep:      time.Millisecond,
	}.NewReadWriter(conn)
	return packeter.Run(prefix.New[uint32](), rw.Pipe)
}

func NewCLIContext(conn io.ReadWriter, parser reflector.Parser[string]) cli.Context {
	pipe := ConnPipe(conn)
	w := channel.Writer{pipe.Snd}

	return cli.NewContext(w, pipe.Rcv, parser)
}

func CLISocket(addr string, rnr cli.CLIRunner) *Socket {
	return New(addr, func(conn net.Conn) {
		ctx := NewCLIContext(conn, nil)
		rnr.Cli(ctx, func() {
			conn.Close()
		})
	})
}
