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

// A caller such as fmt.Fprintf reuses its buffer once Write returns, so what is
// received must not change with it.
func TestWriterReusedBuffer(t *testing.T) {
	ch := make(chan []byte, 2)
	w := channel.Writer{
		Ch: ch,
	}

	buf := []byte("first line")
	n, err := w.Write(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(buf), n)

	copy(buf, "other text")
	w.Write(buf)

	assert.Equal(t, []byte("first line"), <-ch)
	assert.Equal(t, []byte("other text"), <-ch)
}
