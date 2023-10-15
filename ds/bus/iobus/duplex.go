package iobus

import "io"

// Duplex combines an io.Reader and an io.Writer into an io.ReadWriter.
type Duplex struct {
	io.Reader
	io.Writer
}
