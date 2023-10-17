package unixsocket

import (
	"net"
	"strings"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/cli"
)

type RunnerInitilizer func(*cli.Runner)

func (s *Socket) Runner(ri RunnerInitilizer) *Socket {
	s.Handler = func(conn net.Conn) {
		rdr := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReader(conn)

		ec := cli.NewExitClose(
			func() { conn.Close() },
			func() { s.Close() },
		)

		r := &cli.Runner{
			ExitClose: ec,
			Timeout:   25,
			InputProc: func(s string) []string {
				return strings.Split(strings.TrimSpace(s), " ")
			},
			Context: cli.NewContext(conn, rdr.Out, nil),
		}
		ri(r)
		r.Run()
	}
	return s
}
