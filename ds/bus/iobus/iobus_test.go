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

func (b *bufMux) Read(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	b.Lock()
	defer b.Unlock()
	n, err := b.Buffer.Read(p)
	if err == io.EOF {
		err = nil
	}
	return n, err
}

func (b *bufMux) Write(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Write(p)
}

func TestBasic(t *testing.T) {
	buf := &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
	rw := iobus.Config{
		CloseOnEOF: true,
	}.NewReadWriter(buf)

	done := make(chan bool, 1)

	assert.NoError(t, timeout.After(30, func() {
		expected := []byte{1, 2, 3, 4}
		rw.Snd <- expected
		assert.Equal(t, expected, <-rw.Rcv)
	}))

	assert.NoError(t, timeout.After(30, func() {
		expected := []byte{3, 1, 4, 1, 5, 9}
		rw.Snd <- expected
		assert.Equal(t, expected, <-rw.Rcv)
	}))

	assert.NoError(t, timeout.After(200, func() {
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

		done <- true
	}))

	select {
	case err := <-rw.Err:
		assert.NoError(t, err)
	case <-done:
	}

	buf.err = io.EOF

	select {
	case err := <-rw.Err:
		assert.NoError(t, err)
	case in := <-rw.Rcv:
		assert.Nil(t, in)
	case <-time.After(time.Millisecond):
		panic("should get nil from closed in")
	}
}

func TestReadError(t *testing.T) {
	buf := &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
	rw := iobus.NewReadWriter(buf)
	err := fmt.Errorf("this is an error")
	buf.err = err

	assert.Equal(t, err, <-rw.Err)
}

func TestWriteError(t *testing.T) {
	buf := &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
	w, errCh := iobus.NewWriter(buf)
	err := fmt.Errorf("this is an error")
	buf.err = err

	w <- []byte("this is a test")

	assert.Equal(t, err, <-errCh)
}

func TestLong(t *testing.T) {
	msg := make([]byte, 1000)
	rand.Read(msg)

	cfg := iobus.Config{
		Sleep: time.Millisecond,
	}

	buf := bytes.NewBuffer(nil)
	wIn, _ := cfg.NewWriter(buf)

	// by default, message comes through in 2 pieces
	// the first is the size of the buffer
	wIn <- msg
	r := cfg.NewReader(buf)
	timeout.Must(100, func() { assert.Equal(t, msg[:iobus.BufSize], <-r.Out) })
	timeout.Must(100, func() { assert.Equal(t, msg[iobus.BufSize:], <-r.Out) })

	// Using a packeter can reassemble the pieces into the original message
	rwPipe, _, _ := channel.NewPipe(r.Out, wIn)
	prefixPipe := packeter.Run(prefix.New[uint32](), rwPipe)

	prefixPipe.Snd <- msg
	timeout.Must(100, func() { assert.Equal(t, msg, <-prefixPipe.Rcv) })

	r.Stop = true
	close(prefixPipe.Snd)
	timeout.Must(100, func() {
		// check that out is closed
		for range prefixPipe.Rcv {
		}
	})
}
