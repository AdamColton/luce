// Package json32 connects a channel of bytes to a serialbus: the values that are
// sent are serialized as json, with a type32 header to say what they are.
package json32

import (
	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listener"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
)

// Sender combines a type32.TypeMap and a serialbus.Sender. It is setup to send
// json serialized data with a type32 header.
type Sender struct {
	type32.TypeMap
	*serialbus.Sender
}

// NewSender that will write to 'out'. It is setup to use json.Serialize as the
// serializer.
func NewSender(out chan<- []byte) *Sender {
	tm := type32.NewTypeMap()
	return &Sender{
		TypeMap: tm,
		Sender: &serialbus.Sender{
			TypeSerializer: tm.WriterSerializer(json.Serialize),
			Chan:           out,
		},
	}
}

// Receiver combines a serialbus.Receiver with a type32.TypeMap. It is setup to
// receive json serialized data with a type32 header.
type Receiver struct {
	*serialbus.Receiver
	TypeMap type32.TypeMap
	Out     <-chan any
}

// NewReceiver on the 'in' channel. It is setup to receive json serialized data
// with a type32 header.
func NewReceiver(in <-chan []byte) *Receiver {
	iCh := make(chan any)
	tm := type32.NewTypeMap()
	return &Receiver{
		TypeMap: tm,
		Receiver: &serialbus.Receiver{
			In:               in,
			Out:              iCh,
			TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
			TypeRegistrar:    tm,
		},
		Out: iCh,
	}
}

// DefaultMuxSize is the number of handlers that NewHandler makes room for.
const DefaultMuxSize = 10

// NewHandler reads messages off the 'in' channel and sends them to the handler.
// The handler's methods that end in "Handler" handle the types they take, and
// those types are registered with the Receiver. It makes room for
// DefaultMuxSize handlers. It returns the error from registering the handler.
func NewHandler(in <-chan []byte, handler any) (bus.Listener, error) {
	return NewHandlerSize(in, handler, DefaultMuxSize)
}

// NewHandlerSize is NewHandler with room for muxSize handlers.
func NewHandlerSize(in <-chan []byte, handler any, muxSize int) (bus.Listener, error) {
	l, err := listener.New(muxSize, NewReceiver(in), nil)
	if err == nil {
		err = bus.DefaultRegistrar.Register(l, handler)
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}
