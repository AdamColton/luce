package luceio

import (
	"bytes"
	"io"
)

// BufferPool Gets and Puts Buffers, presumably backed by sync.Pool
type BufferPool interface {
	Get() *bytes.Buffer
	Put(buf *bytes.Buffer)
}

// Pool is the default BufferPool.
var Pool BufferPool

func get() *bytes.Buffer {
	if Pool != nil {
		return Pool.Get()
	}
	return &bytes.Buffer{}
}

func put(buf *bytes.Buffer) {
	if Pool != nil {
		Pool.Put(buf)
	}
}

// WriteTo a buffer and return the bytes. If a buffer pool is setup, that will
// be used.
func WriteTo(w io.WriterTo) ([]byte, error) {
	buf := get()
	_, err := w.WriteTo(buf)
	if Pool != nil {
		var b []byte
		if err == nil {
			b := make([]byte, buf.Len())
			copy(b, buf.Bytes())
		}
		put(buf)
		return b, err
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// BufferCloser wraps buffer adding a Close method.
type BufferCloser struct {
	*bytes.Buffer
}

// Close allows BufferCloser to fill the closer interface
func (bc BufferCloser) Close() error {
	return nil
}
