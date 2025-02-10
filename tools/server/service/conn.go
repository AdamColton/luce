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

// Conn wraps a net.Conn with the bus listener and sender used to exchange
// Requests and Responses over it.
type Conn struct {
	Listener bus.Listener
	Sender   *serialbus.Sender
	NetConn  net.Conn
}

// == projects.Code.luce.server ==
// [ ] in-process Conn that skips the unix socket
//  Client/Mux and serviceConn only depend on bus.Sender and bus.Listener, not
//  on NewConn's socket, gob and packeter framing. An in-process service could
//  reuse the same Client/Mux and serviceConn.registerService code paths by
//  wiring them through a trivial, non-serializing bus.Receiver (Run just
//  forwards a chan any to SetOut, RegisterType is a no-op) built with
//  ds/bus/listener.New, and a Sender that pushes straight onto a chan any
//  instead of serializing. That cuts out gob encode/decode and the socket
//  read/write for a service linked into the same binary, leaving only the
//  copy from *http.Request into the Request struct (routeToRequestConverter
//  in tools/server/service.go). serviceConn.Sender in
//  tools/server/servicesocket.go is currently the concrete *serialbus.Sender
//  and would need to widen to the bus.Sender interface to accept either kind
//  of Conn.

// NewConn wraps conn in a gob-serialized, length-prefixed Conn.
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
