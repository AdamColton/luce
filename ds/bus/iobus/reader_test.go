package iobus_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// errReader is an io.Reader that always fails.
type errReader struct {
	err error
}

func (e errReader) Read([]byte) (int, error) {
	return 0, e.err
}

// drain reads from out until it is closed.
func drain(t *testing.T, out <-chan []byte) [][]byte {
	t.Helper()
	var got [][]byte
	assert.NoError(t, timeout.After(1000, func() {
		for b := range out {
			got = append(got, b)
		}
	}))
	return got
}

func TestReader(t *testing.T) {
	buf := newBufMux()
	r := iobus.NewReader(buf)

	expected := []byte("testing")
	buf.Write(expected)
	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, expected, <-r.Out)
	}))
}

func TestReaderStop(t *testing.T) {
	buf := newBufMux()
	r := iobus.Config{
		Sleep:      time.Millisecond,
		CloseOnEOF: false,
	}.NewReader(buf)

	expected := []byte("testing")
	buf.Write(expected)
	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, expected, <-r.Out)
	}))

	r.Stop = true
	assert.Empty(t, drain(t, r.Out), "r.Out should close")
}

func TestReaderCloseOnEOF(t *testing.T) {
	r := iobus.Config{
		CloseOnEOF: true,
	}.NewReader(strings.NewReader("hello"))

	// The data is sent and then r.Out closes.
	assert.Equal(t, [][]byte{[]byte("hello")}, drain(t, r.Out))
}

func TestReaderBufSize(t *testing.T) {
	r := iobus.Config{
		BufSize:    3,
		CloseOnEOF: true,
	}.NewReader(strings.NewReader("abcdefg"))

	expected := [][]byte{[]byte("abc"), []byte("def"), []byte("g")}
	assert.Equal(t, expected, drain(t, r.Out))
}

func TestReaderError(t *testing.T) {
	err := fmt.Errorf("this is an error")
	r := iobus.NewReader(errReader{err})

	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, err, <-r.Err)
	}))
	// After an error the Reader stops and closes r.Out.
	assert.Empty(t, drain(t, r.Out))

	// Without an error channel the Reader still stops.
	r = iobus.Config{}.NewReader(errReader{err})
	assert.Nil(t, r.Err)
	assert.Empty(t, drain(t, r.Out))
}

// Config.Reader can be used directly, with no error channel and no stop flag.
func TestConfigReader(t *testing.T) {
	ch := make(chan []byte)
	go iobus.Config{CloseOnEOF: true}.Reader(strings.NewReader("hi"), ch, nil, nil)

	assert.Equal(t, [][]byte{[]byte("hi")}, drain(t, ch))
}
