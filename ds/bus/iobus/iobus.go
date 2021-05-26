// Package iobus converts an io.Reader or io.Writer to a chan []byte.
package iobus

import (
	"io"
)

// BusReaderConfig sets operational details for a BusReader.
type BusReaderConfig struct {
	CloseOnEOF bool
	BufSize    int
}

// BusReader will read chunks up to size BufSize from io.Reader and write them
// to the In channel
type BusReader struct {
	BusReaderConfig
	In    <-chan []byte
	Err   <-chan error
	in    chan []byte
	errCh chan error
}

// BufSize is the default buffer size that will be used if BufSize is not
// explicitly set on a BusReader.
var BufSize uint = 512

// NewBusReader creates a default BusReader from the provided io.Reader.
func NewBusReader(r io.Reader) *BusReader {
	return BusReaderConfig{
		BufSize: int(BufSize),
	}.New(r)
}

// New creates a BusReader from the provided io.Reader.
func (brc BusReaderConfig) New(r io.Reader) *BusReader {
	br := &BusReader{
		in:    make(chan []byte),
		errCh: make(chan error, 1),
	}
	br.In = br.in
	br.Err = br.errCh
	go br.run(r)

	return br
}

func (br *BusReader) run(r io.Reader) {
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

func (br *BusReader) check(err error) (send, exit bool) {
	if err == nil {
		return true, false
	}
	if err == io.EOF {
		return true, br.CloseOnEOF
	}
	br.errCh <- err
	return false, true
}

// NewBusWriter will write anything sent to the returned channel to the provided
// writer.
func NewBusWriter(w io.Writer) (chan<- []byte, <-chan error) {
	out := make(chan []byte)
	errCh := make(chan error)

	go BusWriter(w, out, errCh)
	return out, errCh
}

// BusWriter reads from a channel and writes anything received to the Writer.
func BusWriter(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	for b := range ch {
		_, err := w.Write(b)
		if err != nil {
			errCh <- err
		}
	}
}

// BusReadWriter runs both a BusReader and a BusWriter on an io.ReaderWriter.
type BusReadWriter struct {
	*BusReader
	Out chan<- []byte
}

// NewBusReadWriter runs both a BusReader and a BusWriter on an io.ReaderWriter.
func NewBusReadWriter(rw io.ReadWriter) *BusReadWriter {
	br := NewBusReader(rw)
	out := make(chan []byte)
	go BusWriter(rw, out, br.errCh)
	return &BusReadWriter{
		BusReader: br,
		Out:       out,
	}
}
