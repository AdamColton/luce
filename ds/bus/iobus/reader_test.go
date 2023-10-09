package iobus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	buf := newBufMux()
	ch, _ := iobus.NewReader(buf)

	expected := []byte("testing")
	buf.Write(expected)
	timeout.After(5, func() {
		assert.Equal(t, expected, <-ch)
	})
}
