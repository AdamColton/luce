package lhttp

import (
	"github.com/adamcolton/luce/math/cmpr"
)

// Socket wraps a MessageReaderWriter, usually a websocket, as an io.Reader and
// io.Writer. Each Write is sent as one text message; Read returns the bytes of
// the messages received, one message after another.
type Socket struct {
	MessageReaderWriter
	buf []byte
}

// NewSocket creates a Socket. It is intended to be used with a websocket.
func NewSocket(socket MessageReaderWriter) *Socket {
	return &Socket{
		MessageReaderWriter: socket,
	}
}

// Read fills p from the current message, reading the next message when the
// current one is used up. An error from ReadMessage is returned as is.
func (socket *Socket) Read(p []byte) (n int, err error) {
	if len(socket.buf) == 0 {
		_, socket.buf, err = socket.ReadMessage()
		if err != nil {
			return 0, err
		}
	}

	ln := cmpr.Min(len(socket.buf), len(p))
	copy(p, socket.buf[:ln])
	socket.buf = socket.buf[ln:]
	return ln, nil
}

// Write sends p as one text message. It returns len(p), or 0 and the error from
// WriteMessage.
func (socket *Socket) Write(p []byte) (n int, err error) {
	err = socket.WriteMessage(1, p)
	if err == nil {
		n = len(p)
	}
	return
}
