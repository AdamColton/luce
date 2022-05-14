package lhttp_test

import (
	"io"
	"testing"

	"github.com/adamcolton/luce/lhttp"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type fakeConn struct {
	in        [][]byte // messages ReadMessage returns, in order
	err       error    // ReadMessage returns this once in is empty, io.EOF if nil
	out       [][]byte // messages written
	failAfter int      // WriteMessage fails once this many were written, if > 0
}

func (f *fakeConn) ReadMessage() (int, []byte, error) {
	if len(f.in) == 0 {
		if f.err == nil {
			return 0, nil, io.EOF
		}
		return 0, nil, f.err
	}
	msg := f.in[0]
	f.in = f.in[1:]
	return websocket.TextMessage, msg, nil
}

func (f *fakeConn) WriteMessage(messageType int, data []byte) error {
	if f.failAfter > 0 && len(f.out) >= f.failAfter {
		return io.ErrClosedPipe
	}
	f.out = append(f.out, data)
	return nil
}

func TestSocketRunReader(t *testing.T) {
	conn := &fakeConn{in: [][]byte{[]byte("hello"), []byte("world")}}
	from := make(chan []byte, 4)
	lhttp.NewSocket(conn).RunReader(from)

	assert.Equal(t, "hello", string(<-from))
	assert.Equal(t, "world", string(<-from))
	_, open := <-from
	assert.False(t, open, "the channel is closed when the socket fails")
}

func TestSocketRunSender(t *testing.T) {
	to := make(chan []byte, 4)
	for _, msg := range []string{"a", "b", "c"} {
		to <- []byte(msg)
	}
	close(to)

	conn := &fakeConn{failAfter: 2}
	lhttp.NewSocket(conn).RunSender(to)
	assert.Equal(t, [][]byte{[]byte("a"), []byte("b")}, conn.out)
	assert.Empty(t, to, "the message that failed was taken off the channel")
}
