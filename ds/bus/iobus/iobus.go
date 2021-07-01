// Package iobus converts an io.Reader or io.Writer to a chan []byte.
package iobus

import (
	"io"
)

// ReaderConfig sets operational details for a Reader.
type ReaderConfig struct {
	CloseOnEOF bool
	BufSize    int
}

// Reader will read chunks up to size BufSize from io.Reader and write them
// to the In channel
type Reader struct {
	ReaderConfig
	In    <-chan []byte
	Err   <-chan error
	in    chan []byte
	errCh chan error
}

// BufSize is the default buffer size that will be used if BufSize is not
// explicitly set on a Reader.
var BufSize uint = 512

// NewReader creates a default Reader from the provided io.Reader.
func NewReader(r io.Reader) *Reader {
	return ReaderConfig{
		BufSize: int(BufSize),
	}.New(r)
}

// New creates a Reader from the provided io.Reader.
func (brc ReaderConfig) New(r io.Reader) *Reader {
	br := &Reader{
		in:           make(chan []byte),
		errCh:        make(chan error, 1),
		ReaderConfig: brc,
	}
	br.In = br.in
	br.Err = br.errCh
	go br.run(r)

	return br
}

func (br *Reader) run(r io.Reader) {
	bufSize := br.BufSize
	if bufSize < 1 {
		bufSize = int(BufSize)
	}

	buf := make([]byte, bufSize)
	for {
		n, err := r.Read(buf)

		send, exit := br.check(err)

		if send && n > 0 {
			br.in <- buf[:n]
			buf = make([]byte, bufSize)
		}

		if exit {
			close(br.in)
			return
		}
	}
}

func (br *Reader) check(err error) (send, exit bool) {
	if err == nil {
		return true, false
	}
	if err == io.EOF {
		return true, br.CloseOnEOF
	}
	br.errCh <- err
	return false, true
}

// NewWriter will write anything sent to the returned channel to the provided
// writer.
func NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	out := make(chan []byte)
	errCh := make(chan error)

	go Writer(w, out, errCh)
	return out, errCh
}

// Writer reads from a channel and writes anything received to the Writer.
func Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	for b := range ch {
		_, err := w.Write(b)
		if err != nil && errCh != nil {
			errCh <- err
		}
	}
}

// ReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter.
type ReadWriter struct {
	*Reader
	Out chan<- []byte
}

// NewReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter.
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	br := NewReader(rw)
	out := make(chan []byte)
	go Writer(rw, out, br.errCh)
	return &ReadWriter{
		Reader: br,
		Out:    out,
	}
}

// NewReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter using
// the ReaderConfig for the reader.
func (brc ReaderConfig) NewReadWriter(rw io.ReadWriter) *ReadWriter {
	br := brc.New(rw)

	out := make(chan []byte)
	go Writer(rw, out, br.errCh)
	return &ReadWriter{
		Reader: br,
		Out:    out,
	}
}
