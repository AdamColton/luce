package iobus

import "time"

// BufSize is the default buffer size that will be used if BufSize is not
// explicitly set on a Reader.
var BufSize uint = 512

// Config sets operational details for a Bus.
type Config struct {
	// CloseOnEOF causes a reader to close when it receives and EOF
	CloseOnEOF bool
	// Buffer size to use
	BufSize int
	// PrefixMessageLength will send 4 bytes indicating the length of the
	// message.
	PrefixMessageLength bool
	// MakeErrCh is used by NewReader and NewWriter to set if an error channel
	// should be created.
	MakeErrCh bool
	// Sleep determines how long to wait before reading again after an empty
	// read.
	Sleep time.Duration
}

func (cfg Config) makeErrCh() chan error {
	if cfg.MakeErrCh {
		return make(chan error, 1)
	}
	return nil
}
