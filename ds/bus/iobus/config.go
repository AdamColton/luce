package iobus

import "time"

// BufSize is the default buffer size, used when Config.BufSize is less than 1.
var BufSize uint = 512

// Config sets operational details for a Reader or Writer.
type Config struct {
	// CloseOnEOF causes a reader to close when it receives an EOF. Otherwise it
	// keeps reading, so data that arrives later is still delivered.
	CloseOnEOF bool
	// BufSize is the size of the buffer for each read. If it is less than 1, the
	// package level BufSize is used.
	BufSize int
	// MakeErrCh is used by the New methods to set if an error channel should be
	// created. The channel has a buffer of 1.
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
