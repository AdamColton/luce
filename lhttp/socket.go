package lhttp

import (
	"github.com/adamcolton/luce/math/cmpr"
)

// Socket is intended to be a websocket
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

func (socket *Socket) Read(p []byte) (n int, err error) {
	ln := len(socket.buf)
	if ln == 0 {
		_, socket.buf, err = socket.ReadMessage()
		if err != nil {
			return 0, err
		}
	}

	ln = cmpr.Min(ln, len(p))
	copy(p, socket.buf[:ln])
	socket.buf = socket.buf[ln:]
	return ln, nil
}

func (socket *Socket) Write(p []byte) (n int, err error) {
	err = socket.WriteMessage(1, p)
	if err == nil {
		n = len(p)
	}
	return
}
