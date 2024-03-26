package channel_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	ch := make(chan []byte)
	w := channel.Writer{
		Ch: ch,
	}

	expected := []byte("hello channel writer")
	go w.Write(expected)

	assert.Equal(t, expected, <-ch)
}
