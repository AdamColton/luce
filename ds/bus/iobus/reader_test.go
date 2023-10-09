package iobus_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	buf := newBufMux()
	r := iobus.NewReader(buf)

	expected := []byte("testing")
	buf.Write(expected)
	timeout.After(5, func() {
		assert.Equal(t, expected, <-r.Out)
	})
}

func TestReaderStop(t *testing.T) {
	buf := newBufMux()
	r := iobus.Config{
		Sleep:      time.Millisecond,
		CloseOnEOF: false,
	}.NewReader(buf)

	expected := []byte("testing")
	buf.Write(expected)
	timeout.After(5, func() {
		assert.Equal(t, expected, <-r.Out)
	})

	r.Stop = true
	timeout.Must(20, func() {
		// check that r.Out has closed
		for x := range <-r.Out {
			_ = x
		}
	})
}
