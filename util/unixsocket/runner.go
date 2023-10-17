package unixsocket

import (
	"net"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/cli"
)

// RunnerInitializer sets up the cli.Runner for a connection. It is called with
// a Runner that has its Context and ExitClose, and sets what the connection
// needs, such as the Commands and the Prompt.
type RunnerInitializer func(*cli.Runner)

// Runner makes the Socket run a cli.Runner for each connection, reading lines
// from the connection and writing to it as text. ri is called for each new
// connection, before its Runner starts, and may be nil. The Runner exits, and
// the connection is closed, when the client exits or disconnects, and the Socket is closed too
// if the client asked to close. Runner replaces the Handler and returns the
// Socket, so that it can be chained after New.
func (s *Socket) Runner(ri RunnerInitializer) *Socket {
	s.Handler = func(conn net.Conn) {
		rdr := iobus.Config{
			// a client that goes away without "exit" ends the Runner
			CloseOnEOF: true,
			Sleep:      time.Millisecond,
		}.NewReader(conn)

		ec := cli.NewExitClose(
			func() { conn.Close() },
			func() { s.Close() },
		)

		r := &cli.Runner{
			ExitClose: ec,
			// the same as cli.NewRunner
			Timeout: 25,
			Context: cli.NewContext(conn, rdr.Out, nil),
		}
		if ri != nil {
			ri(r)
		}
		r.Run()
	}
	return s
}
