package unixsocket_test

import (
	"bytes"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

// syncBuf is what the Context writes to, which the test reads while the Client
// is still running.
type syncBuf struct {
	mux sync.Mutex
	buf strings.Builder
}

func (b *syncBuf) Write(p []byte) (int, error) {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuf) String() string {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.buf.String()
}

// waitFor waits until the buffer contains want.
func (b *syncBuf) waitFor(t *testing.T, want string) {
	err := timeout.After(2000, func() {
		for !strings.Contains(b.String(), want) {
			time.Sleep(time.Millisecond)
		}
	})
	assert.NoError(t, err, "%q not found in %q", want, b.String())
}

// terminal is a Context that the test types into.
func terminal() (cli.Context, *syncBuf, chan []byte) {
	out := &syncBuf{}
	in := make(chan []byte)
	return cli.NewContext(out, in, nil), out, in
}

// echo is a Handler for ConnPipe clients. It answers every message with the
// message after "echo:", until the client sends "bye".
func echo(conn net.Conn) {
	defer conn.Close()
	pipe := unixsocket.ConnPipe(conn)
	for m := range pipe.Rcv {
		if string(m) == "bye" {
			return
		}
		pipe.Snd <- append([]byte("echo:"), m...)
	}
}

// Client looks for the sockets, lists them and asks which to use.
func TestClientChoosesSocket(t *testing.T) {
	fs := lfilemock.Parse(map[string]any{
		"a.sock": "",
		"b.txt":  "",
		"tmp":    map[string]any{"b.sock": ""},
	}).Repository()

	assert.Equal(t, unixsocket.ErrNilContext, unixsocket.ClientFS(nil, fs))

	ctx, out, in := terminal()
	// the user picks a number that is not in the list
	go func() { in <- []byte("5") }()
	assert.Equal(t, unixsocket.ErrNoSuchSocket, unixsocket.ClientFS(ctx, fs))
	assert.Equal(t, "  Sockets:\n    0\ta.sock\n    1\t/tmp/b.sock\n(socket) ", out.String())

	ctx, _, in = terminal()
	go func() { in <- []byte("-1") }()
	assert.Equal(t, unixsocket.ErrNoSuchSocket, unixsocket.ClientFS(ctx, fs))

	// cancelling, with ctrl+x, or closing the input is not an error
	ctx, _, in = terminal()
	go func() { in <- []byte{24} }()
	assert.NoError(t, unixsocket.ClientFS(ctx, fs))
	ctx, _, in = terminal()
	close(in)
	assert.NoError(t, unixsocket.ClientFS(ctx, fs))

	// a socket that nothing is listening at
	ctx, out, in = terminal()
	go func() { in <- []byte("0") }()
	assert.Error(t, unixsocket.ClientFS(ctx, fs))
	assert.Contains(t, out.String(), "  Connecting to a.sock\n\n")
}

// With no sockets there is nothing to ask. Client looks on the real file
// system, so what it finds can't be checked, but it can be cancelled.
func TestClientNoSockets(t *testing.T) {
	assert.Equal(t, unixsocket.ErrNilContext, unixsocket.Client(nil))

	ctx, out, _ := terminal()
	assert.NoError(t, unixsocket.ClientFS(ctx, lfilemock.Parse(nil).Repository()))
	assert.Equal(t, "No sockets found\n", out.String())

	ctx, _, in := terminal()
	close(in)
	assert.NoError(t, unixsocket.Client(ctx))
}

// connectClient serves handler on a socket in a new current directory and
// starts a Client that chooses it. It returns what the Client writes, what it
// reads and the error it returns.
func connectClient(t *testing.T, handler func(net.Conn)) (*syncBuf, chan []byte, <-chan error) {
	// the socket is in the current directory, so its address is short
	t.Chdir(t.TempDir())
	s := unixsocket.New("x.sock", handler)
	errCh := run(s)
	t.Cleanup(func() {
		s.Close()
		<-errCh
	})
	dial(t, "x.sock").Close()

	fs := lfilemock.Parse(map[string]any{"x.sock": ""}).Repository()
	ctx, out, in := terminal()
	done := make(chan error, 1)
	go func() { done <- unixsocket.ClientFS(ctx, fs) }()

	out.waitFor(t, "(socket) ")
	in <- []byte("0")
	out.waitFor(t, "  Connecting to x.sock\n\n")
	return out, in, done
}

// holder is a Handler that sends "ready", does not read, and closes the
// connection when release is closed.
func holder(release <-chan struct{}) func(net.Conn) {
	return func(conn net.Conn) {
		defer conn.Close()
		unixsocket.ConnPipe(conn).Snd <- []byte("ready")
		<-release
	}
}

// awaitReturn waits for the Client to return and checks it has no error.
func awaitReturn(t *testing.T, done <-chan error) {
	assert.NoError(t, timeout.After(2000, func() {
		assert.NoError(t, <-done)
	}))
}

// A Client sends the lines that the user types to the socket, writes what it
// answers, and returns when the socket closes.
func TestClientSession(t *testing.T) {
	out, in, done := connectClient(t, echo)
	in <- []byte("hello")
	out.waitFor(t, "echo:hello")
	in <- []byte("again")
	out.waitFor(t, "echo:again")

	in <- []byte("bye")
	awaitReturn(t, done)
}

// When the input is closed the Client stops reading, and keeps writing what the
// socket sends until the socket closes.
func TestClientInputClosed(t *testing.T) {
	release := make(chan struct{})
	out, in, done := connectClient(t, holder(release))
	close(in)
	out.waitFor(t, "ready")
	close(release)
	awaitReturn(t, done)
}

// A Client returns when the socket closes even if it is sending a line that the
// socket is not taking: the socket is not reading, so the first line fills the
// connection, the second waits in the stages that write it, and the third is
// waiting to be sent.
func TestClientSocketClosesWhileSending(t *testing.T) {
	release := make(chan struct{})
	out, in, done := connectClient(t, holder(release))
	out.waitFor(t, "ready")

	line := bytes.Repeat([]byte("x"), 1<<20)
	for i := 0; i < 3; i++ {
		in <- line
	}
	// the third line is read, and is waiting to be sent
	time.Sleep(20 * time.Millisecond)
	close(release)
	awaitReturn(t, done)
}
