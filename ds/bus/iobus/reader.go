package iobus

import (
	"io"
	"time"
)

// NewReader creates a default Reader from the provided io.Reader.
func NewReader(r io.Reader) (<-chan []byte, <-chan error) {
	return Config{
		MakeErrCh: true,
	}.NewReader(r)
}

// New creates a Reader from the provided io.Reader.
func (cfg Config) NewReader(r io.Reader) (<-chan []byte, <-chan error) {
	ch := make(chan []byte)
	errCh := cfg.makeErrCh()
	go cfg.Reader(r, ch, errCh)

	return ch, errCh
}

// Reader runs a loop reading from r and writing the results to ch.
func (cfg Config) Reader(r io.Reader, ch chan<- []byte, errCh chan<- error) {
	bufSize := cfg.BufSize
	if bufSize < 1 {
		if cfg.PrefixMessageLength {
			bufSize = 4
		} else {
			bufSize = int(BufSize)
		}
	}

	check := makeChecker(cfg.CloseOnEOF, errCh)
	buf := make([]byte, bufSize)

	read := func() (int, error) {
		n, err := r.Read(buf)
		if n == 0 && err == nil {
			time.Sleep(cfg.Sleep)
		}
		return n, err
	}

	for {
		n, err := read()
		send, exit := check(n, err)

		// TODO: replace this with a packeter
		if send && n > 0 && cfg.PrefixMessageLength {
			var ln int
			for i := 0; i < 4 && i < n; i++ {
				ln <<= 8
				ln += int(buf[i])
			}
			buf = make([]byte, ln)
			n, err = read()
			send, exit = check(n, err)
		}

		if send && n > 0 {
			ch <- buf[:n]
			buf = make([]byte, bufSize)
		}

		if exit {
			close(ch)
			return
		}
	}
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
