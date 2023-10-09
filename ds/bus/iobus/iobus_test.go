package iobus_test

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// bufMux is a buffer that different Go routines can read and write. Reading
// from an empty buffer is not an error, it returns nothing. Once setErr is
// called, every Read and Write fails with that error.
type bufMux struct {
	*bytes.Buffer
	sync.Mutex
	err error
}

func newBufMux() *bufMux {
	return &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (b *bufMux) setErr(err error) {
	b.Lock()
	defer b.Unlock()
	b.err = err
}

func (b *bufMux) Read(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	n, err := b.Buffer.Read(p)
	if err == io.EOF {
		err = nil
	}
	return n, err
}

func (b *bufMux) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	return b.Buffer.Write(p)
}

func TestBasic(t *testing.T) {
	buf := newBufMux()
	rw := iobus.Config{
		CloseOnEOF: true,
	}.NewReadWriter(buf)

	// Whatever is sent is written to buf, and then read back.
	assert.NoError(t, timeout.After(1000, func() {
		expected := []byte{1, 2, 3, 4}
		rw.Snd <- expected
		assert.Equal(t, expected, <-rw.Rcv)
	}))

	assert.NoError(t, timeout.After(1000, func() {
		expected := []byte{3, 1, 4, 1, 5, 9}
		rw.Snd <- expected
		assert.Equal(t, expected, <-rw.Rcv)
	}))

	// A message longer than the buffer comes back in pieces the size of the
	// buffer.
	assert.NoError(t, timeout.After(2000, func() {
		expected := make([]byte, 2000)
		rand.Read(expected)
		rw.Snd <- expected
		size := int(iobus.BufSize)
		for i := 0; i < 2000; i += size {
			end := i + size
			if end > 2000 {
				end = 2000
			}
			assert.Equal(t, expected[i:end], <-rw.Rcv)
		}
	}))

	// At EOF the reader closes rw.Rcv.
	buf.setErr(io.EOF)
	assert.NoError(t, timeout.After(1000, func() {
		assert.Nil(t, <-rw.Rcv)
	}))
}

func TestReadError(t *testing.T) {
	buf := newBufMux()
	rw := iobus.NewReadWriter(buf)
	err := fmt.Errorf("this is an error")
	buf.setErr(err)

	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, err, <-rw.Err)
	}))
}

func TestLong(t *testing.T) {
	msg := make([]byte, 1000)
	rand.Read(msg)

	cfg := iobus.Config{
		Sleep: time.Millisecond,
	}

	buf := newBufMux()
	wIn, _ := cfg.NewWriter(buf)

	// by default, message comes through in 2 pieces
	// the first is the size of the buffer
	wIn <- msg
	r := cfg.NewReader(buf)
	timeout.Must(1000, func() { assert.Equal(t, msg[:iobus.BufSize], <-r.Out) })
	timeout.Must(1000, func() { assert.Equal(t, msg[iobus.BufSize:], <-r.Out) })

	// Using a packeter can reassemble the pieces into the original message
	rwPipe, _, _ := channel.NewPipe(r.Out, wIn)
	prefixPipe := packeter.Run(prefix.New[uint32](), rwPipe)

	prefixPipe.Snd <- msg
	timeout.Must(1000, func() { assert.Equal(t, msg, <-prefixPipe.Rcv) })

	r.Stop = true
	close(prefixPipe.Snd)
	timeout.Must(1000, func() {
		// check that out is closed
		for range prefixPipe.Rcv {
		}
	})
}
