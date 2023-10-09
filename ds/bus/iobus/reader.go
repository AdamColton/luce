package iobus

import (
	"io"
	"time"
)

type Reader struct {
	Out  <-chan []byte
	Err  <-chan error
	Stop bool
}

// NewReader creates a default Reader from the provided io.Reader.
func NewReader(r io.Reader) *Reader {
	return Config{
		MakeErrCh: true,
		Sleep:     time.Millisecond,
	}.NewReader(r)
}

// New creates a Reader from the provided io.Reader.
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

// Reader runs a loop reading from r and writing the results to ch.
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
