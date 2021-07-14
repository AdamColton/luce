package iobus

import "io"

// NewReader creates a default Reader from the provided io.Reader.
func NewReader(r io.Reader) (<-chan []byte, <-chan error) {
	return Config{
		BufSize:   int(BufSize),
		MakeErrCh: true,
	}.NewReader(r)
}

// NewWriter will write anything sent to the returned channel to the provided
// writer.
func NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	return Config{
		MakeErrCh: true,
	}.NewWriter(w)
}

// Writer reads from a channel and writes anything received to the Writer.
// If there is an error writing, that will be sent on the errCh, if errCh is
// not nil.
func Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	Config{}.Writer(w, ch, errCh)
}

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
