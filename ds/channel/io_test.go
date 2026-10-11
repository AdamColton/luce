package channel_test

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

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

func TestWriterReadFrom(t *testing.T) {
	ch := make(chan []byte)
	w := channel.Writer{
		Ch: ch,
	}

	// A source longer than one chunk arrives as several messages. Hiding the
	// strings.Reader's WriteTo makes io.Copy use Writer.ReadFrom.
	expected := strings.Repeat("hello channel writer ", 5000)
	var n int64
	var err error
	go func() {
		n, err = io.Copy(w, struct{ io.Reader }{strings.NewReader(expected)})
		close(ch)
	}()

	var got []byte
	for msg := range ch {
		got = append(got, msg...)
	}
	assert.NoError(t, err)
	assert.Equal(t, int64(len(expected)), n)
	assert.Equal(t, expected, string(got))
}

func TestWriterReadFromError(t *testing.T) {
	ch := make(chan []byte, 1)
	w := channel.Writer{
		Ch: ch,
	}

	boom := errors.New("boom")
	n, err := w.ReadFrom(io.MultiReader(strings.NewReader("partial"), iotest.ErrReader(boom)))
	assert.Equal(t, boom, err)
	assert.Equal(t, int64(len("partial")), n)
	assert.Equal(t, []byte("partial"), <-ch)
}

func TestReader(t *testing.T) {
	ch := make(chan []byte)
	r := &channel.Reader{
		Ch: ch,
	}

	go func() {
		ch <- []byte("hello ")
		ch <- nil
		ch <- []byte("channel reader")
		close(ch)
	}()

	// A message is read over several calls if p is short.
	p := make([]byte, 4)
	n, err := r.Read(p)
	assert.NoError(t, err)
	assert.Equal(t, "hell", string(p[:n]))

	n, err = r.Read(nil)
	assert.NoError(t, err)
	assert.Zero(t, n)

	got, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "o channel reader", string(got))

	n, err = r.Read(p)
	assert.Zero(t, n)
	assert.Equal(t, io.EOF, err)
}

// iotest checks the rest of the io.Reader contract.
func TestReaderIOTest(t *testing.T) {
	ch := make(chan []byte, 3)
	ch <- []byte("hello ")
	ch <- []byte("channel ")
	ch <- []byte("reader")
	close(ch)

	r := &channel.Reader{
		Ch: ch,
	}
	assert.NoError(t, iotest.TestReader(r, []byte("hello channel reader")))
}

func TestReaderWriteTo(t *testing.T) {
	ch := make(chan []byte, 3)
	ch <- []byte("hello ")
	ch <- []byte("channel ")
	ch <- []byte("reader")
	close(ch)

	r := &channel.Reader{
		Ch: ch,
	}

	// Part of a message that was already read is not written again.
	p := make([]byte, 3)
	r.Read(p)

	// io.Copy uses Reader.WriteTo.
	var sb strings.Builder
	n, err := io.Copy(&sb, r)
	assert.NoError(t, err)
	assert.Equal(t, int64(len("lo channel reader")), n)
	assert.Equal(t, "lo channel reader", sb.String())
}

type failWriter struct{ err error }

func (f failWriter) Write(p []byte) (int, error) { return len(p) - 1, f.err }

func TestReaderWriteToError(t *testing.T) {
	ch := make(chan []byte, 1)
	ch <- []byte("hello")
	r := &channel.Reader{
		Ch: ch,
	}

	boom := errors.New("boom")
	n, err := r.WriteTo(failWriter{boom})
	assert.Equal(t, boom, err)
	assert.Equal(t, int64(4), n)
}
