package bufpool

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferPool(t *testing.T) {
	buf := Get()
	buf.WriteString("this is a test")
	assert.Equal(t, "this is a test", PutStr(buf))
}

func TestPutStr(t *testing.T) {
	buf := Get()
	buf.WriteString("Hello")
	s := PutStr(buf)
	assert.Equal(t, "Hello", s)
	// another process gets the same buffer
	buf.Reset()
	buf.WriteString("Goodbye")
	assert.Equal(t, "Hello", s)
	buf.Reset()
}

func TestPutAndCopy(t *testing.T) {
	buf := Get()
	buf.WriteString("Hello")
	b := PutAndCopy(buf)
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

func TestMustWriterToString(t *testing.T) {
	s := MustWriterToString(writerto{})
	assert.Equal(t, "testing", s)
}
