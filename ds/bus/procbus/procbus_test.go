package procbus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/procbus"
	"github.com/stretchr/testify/assert"
)

func TestDelim(t *testing.T) {
	in := make(chan []byte)
	out := procbus.Delim(in, '\n')

	in <- []byte("this")
	in <- []byte(" is")
	in <- []byte(" a")
	in <- []byte(" test\nFOOOOO")

	assert.Equal(t, []byte("this is a test"), <-out)
	close(in)
	assert.Equal(t, []byte("FOOOOO"), <-out)
}
