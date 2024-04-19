package unixsocket

import (
	"io"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
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
