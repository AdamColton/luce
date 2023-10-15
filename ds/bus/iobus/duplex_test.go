package iobus_test

import (
	"strings"
	"testing"

	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// A Duplex lets a separate reader and writer be used as one io.ReadWriter.
func TestDuplex(t *testing.T) {
	w := make(chanWriter)
	rw := iobus.Config{
		CloseOnEOF: true,
	}.NewReadWriter(iobus.Duplex{
		Reader: strings.NewReader("incoming"),
		Writer: w,
	})

	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, []byte("incoming"), <-rw.Rcv)
	}))

	rw.Snd <- []byte("outgoing")
	assert.Equal(t, []byte("outgoing"), written(t, w))
}
