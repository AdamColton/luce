package iobus

import "io"

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

	buf := make([]byte, bufSize)
	for {
		n, err := r.Read(buf)
		send, exit := cfg.check(n, err, errCh)
		sum := 0

		if send && n > 0 && cfg.PrefixMessageLength {
			var ln int
			for i := 0; i < 4 && i < n; i++ {
				ln <<= 8
				ln += int(buf[i])
			}
			buf = make([]byte, ln)

			for sum < ln && !exit {
				n, err = r.Read(buf[sum:])
				sum += n
				send, exit = cfg.check(n, err, errCh)
			}
		}

		if send && n > 0 {
			ch <- buf[:sum]
			buf = make([]byte, bufSize)
		}

		if exit {
			close(ch)
			return
		}
	}
}

func (cfg Config) check(n int, err error, errCh chan<- error) (send, exit bool) {
	send = n > 0
	if err == io.EOF {
		exit = cfg.CloseOnEOF
	} else if err != nil {
		if errCh != nil {
			errCh <- err
		}
		send, exit = false, true
	}
	return
}
