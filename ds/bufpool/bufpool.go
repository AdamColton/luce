// Package bufpool reuses bytes.Buffers to reduce allocations.
package bufpool

import (
	"bytes"
	"io"
	"sync"

	"github.com/adamcolton/luce/lerr"
)

// BufferPool provides bytes.Buffers from a pool. A Buffer that is Put back is
// reset, and can be handed out again by a later Get. The Buffer must not be used
// after it is Put.
type BufferPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

type bufferPool struct {
	pool *sync.Pool
}

// Get returns an empty Buffer from the pool, or a new one if the pool has none.
func (b *bufferPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

// Put resets buf and returns it to the pool. buf must not be used after it is
// put.
func (b *bufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	b.pool.Put(buf)
}

// Pool is the package instance of BufferPool. The package level functions use
// it, so replacing it changes where they get their buffers.
var Pool BufferPool = &bufferPool{
	pool: &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	},
}

// Get returns an empty Buffer from Pool.
func Get() *bytes.Buffer { return Pool.Get() }

// Put resets buf and returns it to Pool. buf must not be used after it is put.
func Put(buf *bytes.Buffer) { Pool.Put(buf) }

// PutAndCopy returns buf to the pool and returns a copy of its bytes.
func PutAndCopy(buf *bytes.Buffer) []byte {
	bs := buf.Bytes()
	cp := make([]byte, len(bs))
	copy(cp, bs)
	Pool.Put(buf)
	return cp
}

// PutStr returns buf to the pool and returns its contents as a string.
func PutStr(buf *bytes.Buffer) string {
	s := buf.String() // this makes a copy
	Put(buf)
	return s
}

// WriterToString writes w to a buffer from the pool and returns the result as a
// string. If WriteTo fails, the error is returned along with whatever was
// written.
func WriterToString(w io.WriterTo) (string, error) {
	b := Get()
	_, err := w.WriteTo(b)
	return PutStr(b), err
}

// MustWriterToString writes w to a buffer from the pool and returns the result
// as a string. If WriteTo fails, it panics.
func MustWriterToString(w io.WriterTo) string {
	s, err := WriterToString(w)
	lerr.Panic(err)
	return s
}
