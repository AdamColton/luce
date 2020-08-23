package bufpool

import (
	"bytes"
	"io"
	"sync"

	"github.com/adamcolton/luce/lerr"
)

// BufferPool can Get or Put a Buffer to a pool
type BufferPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

type bufferPool struct {
	pool *sync.Pool
}

// Get returns a Buffer from the pool
func (b *bufferPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

// Put returns a buffer from the pool
func (b *bufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}

// Pool is the package instance of BufferPool.
var Pool BufferPool = &bufferPool{
	pool: &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	},
}

// Get returns a Buffer from the pool
func Get() *bytes.Buffer { return Pool.Get() }

// Put returns a buffer from the pool
func Put(buf *bytes.Buffer) { Pool.Put(buf) }

// PutAndCopy returns the buffer to the pool and returns a copy of it's byte
// slice
func PutAndCopy(buf *bytes.Buffer) []byte {
	bs := buf.Bytes()
	cp := make([]byte, len(bs))
	copy(cp, bs)
	Pool.Put(buf)
	return cp
}

// PutStr returns a buffer from the pool and returns it's value as a string
func PutStr(buf *bytes.Buffer) string {
	s := buf.String() // this makes a copy
	Put(buf)
	return s
}

// WriterToString takes a WriterTo, writes it's contents to a buffer and returns
// the value as a string.
func WriterToString(w io.WriterTo) (string, error) {
	b := Get()
	_, err := w.WriteTo(b)
	return PutStr(b), err
}

// MustWriterToString takes a WriterTo, writes it's contents to a buffer and returns
// the value as a string. If there is an error, it panics.
func MustWriterToString(w io.WriterTo) string {
	s, err := WriterToString(w)
	lerr.Panic(err)
	return s
}
