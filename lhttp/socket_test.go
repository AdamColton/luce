package lhttp_test

import (
	"errors"
	"fmt"
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

func TestSocketRead(t *testing.T) {
	conn := &fakeConn{in: [][]byte{[]byte("hello "), []byte("world!")}}
	got, err := io.ReadAll(lhttp.NewSocket(conn))
	assert.NoError(t, err)
	assert.Equal(t, "hello world!", string(got))
}

func TestSocketReadSmallBuffer(t *testing.T) {
	conn := &fakeConn{in: [][]byte{[]byte("hello")}}
	s := lhttp.NewSocket(conn)
	p := make([]byte, 3)

	n, err := s.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, "hel", string(p[:n]), "the first Read returns data")

	n, err = s.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, "lo", string(p[:n]), "the rest of the message")

	_, err = s.Read(p)
	assert.Equal(t, io.EOF, err)
}

func TestSocketReadError(t *testing.T) {
	boom := errors.New("boom")
	s := lhttp.NewSocket(&fakeConn{err: boom})
	n, err := s.Read(make([]byte, 4))
	assert.Equal(t, 0, n)
	assert.Equal(t, boom, err)
}

func TestSocketWrite(t *testing.T) {
	conn := &fakeConn{}
	s := lhttp.NewSocket(conn)
	n, err := fmt.Fprint(s, "hi there")
	assert.NoError(t, err)
	assert.Equal(t, len("hi there"), n)
	assert.Equal(t, [][]byte{[]byte("hi there")}, conn.out)
}

func TestSocketWriteError(t *testing.T) {
	s := lhttp.NewSocket(&fakeConn{failAfter: 1, out: [][]byte{nil}})
	n, err := s.Write([]byte("hi"))
	assert.Equal(t, 0, n)
	assert.Equal(t, io.ErrClosedPipe, err)
}
