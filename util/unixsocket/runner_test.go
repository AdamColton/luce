package unixsocket_test

import (
	"io"
	"net"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

// readUntil reads from conn until what has been read ends with want, and
// returns it. It fails the test if that takes too long or the connection ends.
func readUntil(t *testing.T, conn net.Conn, want string) string {
	var got strings.Builder
	timeout.Must(2000, func() {
		buf := make([]byte, 1)
		for !strings.HasSuffix(got.String(), want) {
			n, err := conn.Read(buf)
			got.Write(buf[:n])
			if err != nil {
				return
			}
		}
	})
	assert.True(t, strings.HasSuffix(got.String(), want), "%q does not end with %q", got.String(), want)
	return got.String()
}

// A Socket with a Runner is a command line: the client types commands and
// reads the output. "exit" ends the client and "close" ends the Socket too.
func TestSocketRunner(t *testing.T) {
	addr := sockAddr(t)
	var inits atomic.Int32
	s := unixsocket.New(addr, nil).Runner(func(r *cli.Runner) {
		inits.Add(1)
		r.Prompt = "> "
		r.StartMessage = "Welcome\n"
		r.Commands = lerr.Must(handler.Cmds([]*handler.Command{
			{
				Name:   "hi",
				Action: func() { r.WriteString("Hi!") },
			}, {
				Name:   "exit",
				Action: func() { r.Exit = true },
			}, {
				Name: "close",
				Action: func() {
					r.Close = true
					r.Exit = true
				},
			},
		}))
	})
	errCh := run(s)

	conn := dial(t, addr)
	readUntil(t, conn, "Welcome\n> ")
	_, err := conn.Write([]byte("hi\n"))
	assert.NoError(t, err)
	readUntil(t, conn, "Hi!\n> ")
	_, err = conn.Write([]byte("exit\n"))
	assert.NoError(t, err)

	// exit closes that connection and leaves the Socket running
	var rest []byte
	timeout.Must(2000, func() {
		rest, err = io.ReadAll(conn)
	})
	assert.NoError(t, err)
	assert.Equal(t, "\n", string(rest))

	conn = dial(t, addr)
	readUntil(t, conn, "Welcome\n> ")
	_, err = conn.Write([]byte("close\n"))
	assert.NoError(t, err)
	assert.NoError(t, timeout.After(2000, func() {
		assert.NoError(t, <-errCh)
	}))
	assert.Equal(t, int32(2), inits.Load())
}

// The initializer is optional. Without one the Runner has no Commands, which
// is enough to see the connection is served until the client goes away.
func TestSocketRunnerNoInit(t *testing.T) {
	client, server := net.Pipe()
	s := unixsocket.New("unused", nil).Runner(nil)
	done := make(chan struct{})
	go func() {
		s.Handler(server)
		close(done)
	}()

	assert.NoError(t, client.Close())
	assert.NoError(t, timeout.After(2000, done))
}
