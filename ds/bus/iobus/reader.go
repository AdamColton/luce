package iobus

import (
	"io"
	"time"
)

// Reader is a Go routine that reads from an io.Reader and sends what it reads
// on Out. Create one with NewReader.
type Reader struct {
	// Out receives the data that was read. Each slice is newly allocated. Out is
	// closed when the Go routine finishes.
	Out <-chan []byte
	// Err receives an error if reading fails. It is nil unless the Config asked
	// for an error channel.
	Err <-chan error
	// Stop ends the loop the next time it checks, which is after each read. It is
	// not synchronized with the Go routine that reads it.
	Stop bool
}

// NewReader creates a Reader from the provided io.Reader. It has an error
// channel and sleeps for a millisecond after an empty read. It does not close at
// EOF.
func NewReader(r io.Reader) *Reader {
	return Config{
		MakeErrCh: true,
		Sleep:     time.Millisecond,
	}.NewReader(r)
}

// NewReader creates a Reader from the provided io.Reader using cfg.
func (cfg Config) NewReader(r io.Reader) *Reader {
	ch := make(chan []byte)
	errCh := cfg.makeErrCh()
	out := &Reader{
		Out: ch,
		Err: errCh,
	}
	go cfg.Reader(r, ch, errCh, &(out.Stop))

	return out
}

// Reader runs a loop reading from r and sending what it reads on ch. The loop
// ends when *stop is true (it is checked after each read), when a read fails, or
// at EOF if CloseOnEOF is set. Then ch is closed. If a read fails the error is
// sent on errCh, if it is not nil. If stop is nil the loop can only end with an
// error or at EOF. The slices sent on ch are not reused.
func (cfg Config) Reader(r io.Reader, ch chan<- []byte, errCh chan<- error, stop *bool) {
	bufSize := cfg.BufSize
	if bufSize < 1 {
		bufSize = int(BufSize)
	}

	if stop == nil {
		s := false
		stop = &s
	}

	check := makeChecker(cfg.CloseOnEOF, errCh)
	buf := make([]byte, bufSize)

	read := func() (int, error) {
		n, err := r.Read(buf)
		return n, err
	}

	for !*stop {
		n, err := read()
		send, exit := check(n, err)
		if send && n > 0 {
			ch <- buf[:n]
			buf = make([]byte, bufSize)
		}

		if exit {
			break
		}

		if n == 0 && cfg.Sleep > 0 {
			time.Sleep(cfg.Sleep)
		}
	}

	close(ch)
}

type checker func(n int, err error) (send, exit bool)

func makeChecker(closeOnEOF bool, errCh chan<- error) checker {
	return func(n int, err error) (send, exit bool) {
		send = n > 0
		if err == io.EOF {
			exit = closeOnEOF
		} else if err != nil {
			if errCh != nil {
				errCh <- err
			}
			send, exit = false, true
		}
		return
	}
}
