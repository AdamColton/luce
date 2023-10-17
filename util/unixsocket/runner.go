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
		in, _ := iobus.Config{
			Sleep: time.Millisecond,
		}.NewReader(conn)

		r := &cli.Runner{
			OnExit:  func() { conn.Close() },
			OnClose: func() { s.Close() },
			Timeout: 25,
			InputProc: func(s string) []string {
				return strings.Split(strings.TrimSpace(s), " ")
			},
			Context: cli.NewContext(conn, in, nil),
		}
		ri(r)
		r.Run()
	}
	return s
}
