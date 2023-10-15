package iobus

import "io"

type Duplex struct {
	io.Reader
	io.Writer
}
