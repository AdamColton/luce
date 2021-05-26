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
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

type bufMux struct {
	*bytes.Buffer
	sync.Mutex
	err error
}

func (b *bufMux) Read(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Read(p)
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
	rw := iobus.NewBusReadWriter(buf)

	done := make(chan bool, 1)

	assert.NoError(t, timeout.After(30, func() {
		expected := []byte{1, 2, 3, 4}
		rw.Out <- expected
		assert.Equal(t, expected, <-rw.In)
	}))

	assert.NoError(t, timeout.After(30, func() {
		expected := []byte{3, 1, 4, 1, 5, 9}
		rw.Out <- expected
		assert.Equal(t, expected, <-rw.In)
	}))

	assert.NoError(t, timeout.After(200, func() {
		expected := make([]byte, 2000)
		rand.Read(expected)
		rw.Out <- expected
		size := int(iobus.BufSize)
		for i := 0; i < 2000; i += size {
			end := i + size
			if end > 2000 {
				end = 2000
			}
			assert.Equal(t, expected[i:end], <-rw.In)
		}

		done <- true
	}))

	select {
	case err := <-rw.Err:
		assert.NoError(t, err)
	case <-done:
	}

	rw.CloseOnEOF = true
	buf.err = io.EOF

	select {
	case err := <-rw.Err:
		assert.NoError(t, err)
	case in := <-rw.In:
		assert.Nil(t, in)
	case <-time.After(time.Millisecond):
		panic("should get nil from closed in")
	}
}

func TestReadError(t *testing.T) {
	buf := &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
	rw := iobus.NewBusReadWriter(buf)
	err := fmt.Errorf("this is an error")
	buf.err = err

	assert.Equal(t, err, <-rw.Err)
}

func TestWriteError(t *testing.T) {
	buf := &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
	w, errCh := iobus.NewBusWriter(buf)
	err := fmt.Errorf("this is an error")
	buf.err = err

	w <- []byte("this is a test")

	assert.Equal(t, err, <-errCh)
}
