package iobus_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// chanWriter is an io.Writer that sends a copy of everything written to it on
// the channel, so a test can wait for a Write.
type chanWriter chan []byte

func (c chanWriter) Write(p []byte) (int, error) {
	c <- append([]byte(nil), p...)
	return len(p), nil
}

// failWriter is an io.Writer that always fails. It reports each attempt on
// tried.
type failWriter struct {
	err   error
	tried chan struct{}
}

func (f failWriter) Write([]byte) (int, error) {
	f.tried <- struct{}{}
	return 0, f.err
}

// written gets what was written to w.
func written(t *testing.T, w chanWriter) []byte {
	t.Helper()
	var got []byte
	assert.NoError(t, timeout.After(1000, func() {
		got = <-w
	}))
	return got
}

func TestWriter(t *testing.T) {
	w := make(chanWriter)
	ch := make(chan []byte)
	done := make(chan struct{})
	go func() {
		iobus.Writer(w, ch, nil)
		close(done)
	}()

	ch <- []byte("one")
	assert.Equal(t, []byte("one"), written(t, w))
	ch <- []byte("two")
	assert.Equal(t, []byte("two"), written(t, w))

	// Writer returns when its channel is closed.
	close(ch)
	assert.NoError(t, timeout.After(1000, done))
}

func TestNewWriter(t *testing.T) {
	w := make(chanWriter)
	in, errCh := iobus.NewWriter(w)
	assert.NotNil(t, errCh)

	in <- []byte("testing")
	assert.Equal(t, []byte("testing"), written(t, w))
	close(in)

	// A Config without MakeErrCh has no error channel.
	in, errCh = iobus.Config{}.NewWriter(w)
	assert.Nil(t, errCh)
	in <- []byte("more")
	assert.Equal(t, []byte("more"), written(t, w))
	close(in)
}

func TestWriterError(t *testing.T) {
	err := fmt.Errorf("this is an error")
	in, errCh := iobus.NewWriter(failWriter{err, make(chan struct{}, 1)})

	in <- []byte("this is a test")
	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, err, <-errCh)
	}))
}

// Without an error channel, a failed write is dropped and the Writer carries on.
func TestWriterErrorNoChannel(t *testing.T) {
	tried := make(chan struct{})
	in, errCh := iobus.Config{}.NewWriter(failWriter{fmt.Errorf("failed"), tried})
	assert.Nil(t, errCh)

	for i := 0; i < 2; i++ {
		in <- []byte("this is a test")
		assert.NoError(t, timeout.After(1000, tried))
	}
	close(in)
}
