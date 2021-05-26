package unixsocket_test

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

// sockAddr is a socket address in a new directory. The directory is under the
// system temp directory and not t.TempDir, because the name of the test can
// make the path longer than a unix socket address can be.
func sockAddr(t *testing.T) string {
	dir, err := os.MkdirTemp("", "us")
	assert.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(dir) })
	return filepath.Join(dir, "test.sock")
}

// dial connects to addr, waiting for the socket to be listening.
func dial(t *testing.T, addr string) (conn net.Conn) {
	timeout.Must(2000, func() {
		for {
			var err error
			if conn, err = net.Dial("unix", addr); err == nil {
				return
			}
			time.Sleep(time.Millisecond)
		}
	})
	t.Cleanup(func() { conn.Close() })
	return
}

// run calls Run in a Go routine, the error is sent on the returned channel.
func run(s *unixsocket.Socket) <-chan error {
	errCh := make(chan error, 1)
	go func() { errCh <- s.Run() }()
	return errCh
}

// upper is a Handler that answers five bytes with the same five in upper case.
func upper(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 5)
	if _, err := io.ReadFull(conn, buf); err == nil {
		conn.Write(bytes.ToUpper(buf))
	}
}

// ask sends a word to the Socket at addr and returns what comes back.
func ask(t *testing.T, addr, word string) string {
	conn := dial(t, addr)
	_, err := conn.Write([]byte(word))
	assert.NoError(t, err)
	buf := make([]byte, len(word))
	timeout.Must(2000, func() {
		_, err = io.ReadFull(conn, buf)
	})
	assert.NoError(t, err)
	return string(buf)
}

// fileAt makes a mock tree with a file at path.
func fileAt(path string) *lfilemock.Directory {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	tree := map[string]any{parts[len(parts)-1]: ""}
	for i := len(parts) - 2; i >= 0; i-- {
		tree = map[string]any{parts[i]: tree}
	}
	return lfilemock.Parse(tree)
}

// A Socket clears the file a previous run left at its address before it
// listens. The addresses are in a directory that does not exist, so Run stops at
// listening, and what happens before that is what is tested.
func TestSocketClearsStaleFile(t *testing.T) {
	addr := filepath.ToSlash(filepath.Join(t.TempDir(), "nodir", "test.sock"))
	handler := func(net.Conn) {}
	listenErr := func(err error) {
		var opErr *net.OpError
		if assert.True(t, errors.As(err, &opErr), "%v", err) {
			assert.Equal(t, "listen", opErr.Op)
		}
	}

	// the file a previous run left is removed
	root := fileAt(addr)
	s := unixsocket.New(addr, handler)
	s.FS = root.Repository()
	_, found := root.Get(addr)
	assert.True(t, found)
	listenErr(s.Run())
	_, found = root.Get(addr)
	assert.False(t, found)

	// there is nothing to remove the second time, which is not an error
	listenErr(s.Run())

	// any other error from the file system is the error of Run
	root.Err = lerr.Str("boom")
	assert.Equal(t, lerr.Str("boom"), s.Run())
}

// A Socket answers connections until it is closed. Closing it waits for Run to
// finish, it removes the socket file, it can be done again and at the same time,
// and the Socket can run again.
func TestSocketRunClose(t *testing.T) {
	addr := sockAddr(t)
	s := unixsocket.New(addr, upper)

	// Close before Run has nothing to do
	s.Close()

	for i := 0; i < 2; i++ {
		errCh := run(s)
		assert.Equal(t, "HELLO", ask(t, addr, "hello"))
		assert.Equal(t, "WORLD", ask(t, addr, "world"))

		// it is already running
		assert.Equal(t, unixsocket.ErrRunning, s.Run())

		var wg sync.WaitGroup
		wg.Add(2)
		for j := 0; j < 2; j++ {
			go func() {
				s.Close()
				wg.Done()
			}()
		}
		assert.NoError(t, timeout.After(2000, &wg))
		assert.NoError(t, <-errCh)

		_, err := os.Stat(addr)
		assert.ErrorIs(t, err, fs.ErrNotExist)
		s.Close()
	}
}

// Run needs a Handler, and says so rather than failing on the first connection.
func TestSocketNilHandler(t *testing.T) {
	addr := sockAddr(t)
	s := unixsocket.New(addr, nil)
	assert.Equal(t, unixsocket.ErrNilHandler, s.Run())

	s.Handler = upper
	errCh := run(s)
	assert.Equal(t, "HELLO", ask(t, addr, "hello"))
	s.Close()
	assert.NoError(t, <-errCh)
}
