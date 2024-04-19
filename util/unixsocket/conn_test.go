package unixsocket_test

import (
	"net"
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

// Two ConnPipes on the two ends of a connection carry messages in both
// directions, and each message arrives whole.
func TestConnPipe(t *testing.T) {
	a, b := net.Pipe()
	left, right := unixsocket.ConnPipe(a), unixsocket.ConnPipe(b)

	timeout.Must(2000, func() {
		left.Snd <- []byte("hello")
		assert.Equal(t, "hello", string(<-right.Rcv))
		right.Snd <- []byte("world, and more")
		assert.Equal(t, "world, and more", string(<-left.Rcv))
	})

	// when the connection ends, Rcv is closed
	assert.NoError(t, b.Close())
	timeout.Must(2000, func() {
		_, ok := <-left.Rcv
		assert.False(t, ok)
	})
}
