package bufpool_test

import (
	"io"
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestBufferPool(t *testing.T) {
	buf := bufpool.Get()
	buf.WriteString("this is a test")
	assert.Equal(t, "this is a test", bufpool.PutStr(buf))
}

func TestPut(t *testing.T) {
	buf := bufpool.Get()
	buf.WriteString("this is a test")
	bufpool.Put(buf)

	// Whatever buffer comes out of the pool is empty.
	buf = bufpool.Get()
	assert.Zero(t, buf.Len())
	bufpool.Put(buf)
}

func TestPutStr(t *testing.T) {
	buf := bufpool.Get()
	buf.WriteString("Hello")
	s := bufpool.PutStr(buf)
	assert.Equal(t, "Hello", s)
	// another process gets the same buffer
	buf.Reset()
	buf.WriteString("Goodbye")
	assert.Equal(t, "Hello", s)
	buf.Reset()
}

func TestPutAndCopy(t *testing.T) {
	buf := bufpool.Get()
	buf.WriteString("Hello")
	b := bufpool.PutAndCopy(buf)
	assert.Equal(t, []byte("Hello"), b)
	// another process gets the same buffer
	buf.Reset()
	buf.WriteString("Goodbye")
	assert.Equal(t, []byte("Hello"), b)
	buf.Reset()
}

type writerto struct{}

func (writerto) WriteTo(w io.Writer) (int64, error) {
	w.Write([]byte("testing"))
	return 7, nil
}

const errFailed = lerr.Str("failed")

type failing struct{}

func (failing) WriteTo(w io.Writer) (int64, error) {
	w.Write([]byte("partial"))
	return 7, errFailed
}

func TestWriterToString(t *testing.T) {
	s, err := bufpool.WriterToString(writerto{})
	assert.NoError(t, err)
	assert.Equal(t, "testing", s)

	// The error is returned along with whatever was written.
	s, err = bufpool.WriterToString(failing{})
	assert.ErrorIs(t, err, errFailed)
	assert.Equal(t, "partial", s)
}

func TestMustWriterToString(t *testing.T) {
	assert.Equal(t, "testing", bufpool.MustWriterToString(writerto{}))
	assert.Panics(t, func() {
		bufpool.MustWriterToString(failing{})
	})
}
