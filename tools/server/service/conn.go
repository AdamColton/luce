package service

import (
	"net"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/wrap/gob"
)

// Conn bundles a bus listener and sender witha NetConn for communicating over
// a Unix socket. It uses type32 prefixing with Gob for serialization.
type Conn struct {
	Listener bus.Listener
	Sender   *serialbus.Sender
	NetConn  net.Conn
}

// NewConn created from the underlying net.Conn. It uses type32 prefixing with
// Gob for serialization.
func NewConn(conn net.Conn) (*Conn, error) {
	b := iobus.Config{
		CloseOnEOF:          true,
		PrefixMessageLength: true,
	}.NewReadWriter(conn)

	l, err := serialbus.NewListener(b.In, tm.ReaderDeserializer(gob.Deserialize), tm, nil)
	if err != nil {
		return nil, err
	}
	s := &serialbus.Sender{
		TypeSerializer: tm.WriterSerializer(gob.Serialize),
		Chan:           b.Out,
	}

	return &Conn{
		Listener: l,
		Sender:   s,
		NetConn:  conn,
	}, nil
}
