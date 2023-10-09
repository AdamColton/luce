package iobus

import (
	"io"

	"github.com/adamcolton/luce/ds/channel"
)

// ReadWriter reads from and writes to an io.ReadWriter using channels: Rcv
// receives what is read and everything sent on Snd is written.
type ReadWriter struct {
	channel.Pipe[[]byte]
	// Err receives errors from both the reading and the writing. It is nil unless
	// the Config asked for an error channel.
	Err <-chan error
	// Stop ends the reading loop, as Reader.Stop does.
	Stop bool
}

// NewReadWriter runs both a Reader and a Writer on an io.ReadWriter. It has an
// error channel. Unlike NewReader it does not sleep after an empty read.
func NewReadWriter(rw io.ReadWriter) *ReadWriter {
	return Config{
		MakeErrCh: true,
	}.NewReadWriter(rw)
}

// NewReadWriter runs both a Reader and a Writer on an io.ReadWriter using cfg.
// They share the error channel.
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
