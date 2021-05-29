package service

import (
	"net"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/iobus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
)

type Conn struct {
	Listener bus.Listener
	Sender   *serialbus.Sender
	NetConn  net.Conn
}

func NewConn(conn net.Conn) (*Conn, error) {
	rw := iobus.Config{
		CloseOnEOF: true,
	}.NewReadWriter(conn)
	p := packeter.Run(prefix.New[uint32](), rw.Pipe)

	l, err := serialbus.NewListener(p.Rcv, tm.ReaderDeserializer(gob.Deserialize), tm, nil)
	if err != nil {
		return nil, err
	}
	s := &serialbus.Sender{
		TypeSerializer: tm.WriterSerializer(gob.Serialize),
		Chan:           p.Snd,
	}

	return &Conn{
		Listener: l,
		Sender:   s,
		NetConn:  conn,
	}, nil
}
