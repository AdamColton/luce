package iobus_test

import (
	"bytes"
	"io"
	"sync"
)

// bufMux is a buffer that different Go routines can read and write. Reading
// from an empty buffer is not an error, it returns nothing. Once setErr is
// called, every Read and Write fails with that error.
type bufMux struct {
	*bytes.Buffer
	sync.Mutex
	err error
}

func newBufMux() *bufMux {
	return &bufMux{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (b *bufMux) setErr(err error) {
	b.Lock()
	defer b.Unlock()
	b.err = err
}

func (b *bufMux) Read(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	n, err := b.Buffer.Read(p)
	if err == io.EOF {
		err = nil
	}
	return n, err
}

func (b *bufMux) Write(p []byte) (int, error) {
	b.Lock()
	defer b.Unlock()
	if b.err != nil {
		return 0, b.err
	}
	return b.Buffer.Write(p)
}
