package json32

import (
	"github.com/adamcolton/luce/ds/bus"
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
	Out     <-chan interface{}
}

// NewReceiver on the 'in' channel. It is setup to receive json serialized data
// with a type32 header.
func NewReceiver(in <-chan []byte) *Receiver {
	iCh := make(chan interface{})
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

// NewHandler reads messages off the 'in' channel and sends them to the handler.
func NewHandler(in <-chan []byte, handler interface{}) (bus.Listener, error) {
	r := NewReceiver(in)
	l, err := bus.NewListener(r, nil, nil)
	if err != nil {
		return nil, err
	}
	err = bus.RegisterHandlerType(l, handler)
	if err != nil {
		return nil, err
	}
	return l, nil
}
