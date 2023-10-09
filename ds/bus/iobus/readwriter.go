package iobus

import "io"

// ReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter.
type ReadWriter struct {
	In  <-chan []byte
	Out chan<- []byte
	Err <-chan error
}

// NewReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter.
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return Config{
		MakeErrCh: true,
	}.NewReadWriter(rw)
}

// NewReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter using
// the ReaderConfig for the reader.
func (cfg Config) NewReadWriter(rw io.ReadWriter) *ReadWriter {
	in := make(chan []byte)
	out := make(chan []byte)
	errCh := cfg.makeErrCh()

	go cfg.Reader(rw, in, errCh)
	go cfg.Writer(rw, out, errCh)

	return &ReadWriter{
		In:  in,
		Out: out,
		Err: errCh,
	}
}
