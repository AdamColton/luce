package unixsocket

import (
	"net"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
)

func ConnPipe(conn net.Conn) channel.Pipe[[]byte] {
	rw := iobus.Config{
		CloseOnEOF: true,
		Sleep:      time.Millisecond,
	}.NewReadWriter(conn)
	return packeter.Run(prefix.New[uint32](), rw.Pipe)
}
