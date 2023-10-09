package iobus_test

import (
	"bytes"
	"io"
	"sync"
)

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

func (b *bufMux) Read(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	b.Lock()
	defer b.Unlock()
	n, err := b.Buffer.Read(p)
	if err == io.EOF {
		err = nil
	}
	return n, err
}

func (b *bufMux) Write(p []byte) (int, error) {
	if b.err != nil {
		return 0, b.err
	}
	b.Lock()
	defer b.Unlock()
	return b.Buffer.Write(p)
}
