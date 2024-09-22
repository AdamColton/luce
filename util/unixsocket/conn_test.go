package unixsocket_test

import (
	"net"
	"testing"

	"github.com/adamcolton/luce/util/cli"
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

// NewCLIContext is a cli.Context on a connection: the messages that arrive are
// the lines it reads, and what it writes is sent as messages.
func TestNewCLIContext(t *testing.T) {
	a, b := net.Pipe()
	ctx := unixsocket.NewCLIContext(a, nil)
	other := unixsocket.ConnPipe(b)

	timeout.Must(2000, func() {
		other.Snd <- []byte("  list files \n")
		assert.Equal(t, "list files", ctx.ReadString(nil))

		go ctx.WriteStrings("a.txt", " b.txt")
		assert.Equal(t, "a.txt", string(<-other.Rcv))
		assert.Equal(t, " b.txt", string(<-other.Rcv))

		// the parser is for Input
		other.Snd <- []byte("42")
		var n int
		go func() { <-other.Rcv }()
		assert.True(t, ctx.Input("n: ", &n))
		assert.Equal(t, 42, n)
	})
	assert.Equal(t, cli.Parser, ctx.Parser(), "nil parser is the default")
}

// cliRunner is a CLIRunner that answers each line with the line repeated, and
// exits on "bye".
type cliRunner struct{}

func (cliRunner) Cli(ctx cli.Context, onExit func()) {
	for {
		line := ctx.ReadString(nil)
		if ctx.Closed() || line == "bye" {
			onExit()
			return
		}
		ctx.WriteString(line + line)
	}
}

// A CLISocket runs the CLIRunner for each connection. A Client talks to it.
func TestCLISocket(t *testing.T) {
	addr := sockAddr(t)
	s := unixsocket.CLISocket(addr, cliRunner{})
	errCh := run(s)
	defer func() {
		s.Close()
		assert.NoError(t, <-errCh)
	}()

	pipe := unixsocket.ConnPipe(dial(t, addr))
	timeout.Must(2000, func() {
		pipe.Snd <- []byte("ab")
		assert.Equal(t, "abab", string(<-pipe.Rcv))
		pipe.Snd <- []byte("bye")
		// the connection is closed when the runner exits
		_, ok := <-pipe.Rcv
		assert.False(t, ok)
	})
}
