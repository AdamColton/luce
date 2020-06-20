package bufpool_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"
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
