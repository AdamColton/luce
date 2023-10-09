package iobus

import (
	"io"
)

// Writer reads from ch and writes anything received to w until ch is closed. If
// there is an error writing, that will be sent on the errCh, if errCh is not
// nil, and Writer carries on with the next value. The send on errCh blocks if
// nothing is receiving from it.
func Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	Config{}.Writer(w, ch, errCh)
}

// NewWriter will write anything sent to the returned channel to the provided
// writer. The returned error channel has a buffer of 1. Closing the returned
// channel ends the Go routine.
func NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	return Config{
		MakeErrCh: true,
	}.NewWriter(w)
}

// NewWriter creates the channels and runs Writer in a Go routine. The error
// channel is nil unless cfg.MakeErrCh is set.
func (cfg Config) NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	out := make(chan []byte)
	errCh := cfg.makeErrCh()

	go cfg.Writer(w, out, errCh)
	return out, errCh
}

// Writer is Writer, run with cfg. The Config does not change how it writes.
func (cfg Config) Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	for b := range ch {
		_, err := w.Write(b)
		if err != nil && errCh != nil {
			errCh <- err
		}
	}
}
