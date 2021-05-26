package unixsocket

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

// failListener is a listener that can't accept.
type failListener struct{ err error }

func (f failListener) Accept() (net.Conn, error) { return nil, f.err }
func (failListener) Close() error                { return nil }
func (failListener) Addr() net.Addr              { return nil }

// removed is a FileSystem that notes what was removed.
type removed []string

func (r *removed) Remove(name string) error {
	*r = append(*r, name)
	return nil
}

// This is internal because a listener that fails to accept can't be made with a
// real socket. When Accept fails for a reason other than Close, Run returns the
// error, it cleans up as it does for Close, and it can run again.
func TestRunAcceptError(t *testing.T) {
	boom := errors.New("boom")
	fs := &removed{}
	s := New("fail.sock", func(net.Conn) {})
	s.FS = fs
	s.listen = func(network, addr string) (net.Listener, error) {
		assert.Equal(t, "unix", network)
		assert.Equal(t, "fail.sock", addr)
		return failListener{boom}, nil
	}

	assert.Equal(t, boom, s.Run())
	// once before it listens and once when it stopped
	assert.Equal(t, removed{"fail.sock", "fail.sock"}, *fs)
	assert.Equal(t, boom, s.Run(), "not ErrRunning, and not stuck")
}
