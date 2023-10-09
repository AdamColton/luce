// Package iobus converts an io.Reader or io.Writer to a chan []byte.
package iobus

import (
	"io"
)

// Writer reads from a channel and writes anything received to the Writer.
// If there is an error writing, that will be sent on the errCh, if errCh is
// not nil.
func Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	Config{}.Writer(w, ch, errCh)
}

// NewWriter will write anything sent to the returned channel to the provided
// writer.
func NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	return Config{
		MakeErrCh: true,
	}.NewWriter(w)
}

// NewWriter creates the channels and runs Writer in a Go routine.
func (cfg Config) NewWriter(w io.Writer) (chan<- []byte, <-chan error) {
	out := make(chan []byte)
	errCh := cfg.makeErrCh()

	go cfg.Writer(w, out, errCh)
	return out, errCh
}

// Writer reads from the channel and writes anything it receives to the writer.
func (cfg Config) Writer(w io.Writer, ch <-chan []byte, errCh chan<- error) {
	for b := range ch {
		if cfg.PrefixMessageLength {
			lnMsg := make([]byte, 4)
			ln := len(b)
			for i := 3; i >= 0 && ln > 0; i-- {
				lnMsg[i] = byte(ln)
				ln >>= 8
			}
			_, err := w.Write(lnMsg)
			if err != nil && errCh != nil {
				errCh <- err
			}
		}
		_, err := w.Write(b)
		if err != nil && errCh != nil {
			errCh <- err
		}
	}
}
