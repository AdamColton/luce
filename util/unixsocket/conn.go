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

// ConnPipe turns a connection into a channel.Pipe of messages. A message is
// sent as a packet with a uint32 length prefix, so the other end must read it
// the same way, and a message that is received is a whole packet. The Rcv channel
// is closed when the connection ends.
func ConnPipe(conn io.ReadWriter) channel.Pipe[[]byte] {
	rw := iobus.Config{
		CloseOnEOF: true,
		Sleep:      time.Millisecond,
	}.NewReadWriter(conn)
	return packeter.Run(prefix.New[uint32](), rw.Pipe)
}

// NewCLIContext creates a cli.Context that reads and writes messages on conn, as
// ConnPipe does. Each write is sent as one message, and each message that
// arrives is a line to read. If parser is nil the Context uses cli.Parser.
func NewCLIContext(conn io.ReadWriter, parser reflector.Parser[string]) cli.Context {
	pipe := ConnPipe(conn)
	w := channel.Writer{pipe.Snd}

	return cli.NewContext(w, pipe.Rcv, parser)
}

// CLISocket creates a Socket that runs rnr on a cli.Context, made with
// NewCLIContext, for each connection. The connection is closed when rnr exits.
// It is the server that Client connects to.
func CLISocket(addr string, rnr cli.CLIRunner) *Socket {
	return New(addr, func(conn net.Conn) {
		ctx := NewCLIContext(conn, nil)
		rnr.Cli(ctx, func() {
			conn.Close()
		})
	})
}
