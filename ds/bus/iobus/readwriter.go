package iobus

import (
	"io"

	"github.com/adamcolton/luce/ds/channel"
)

// ReadWriter runs both a Reader and a BusWriter on an io.ReaderWriter.
type ReadWriter struct {
	channel.Pipe[[]byte]
	Err  <-chan error
	Stop bool
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

	ret := &ReadWriter{
		Pipe: channel.Pipe[[]byte]{
			Rcv: in,
			Snd: out,
		},
		Err: errCh,
	}

	go cfg.Reader(rw, in, errCh, &(ret.Stop))
	go cfg.Writer(rw, out, errCh)

	return ret
}
