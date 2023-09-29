package serialbus

import (
	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listener"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// Receiver takes serialized messages off a serial bus, deserializes them
// and places the deserialized objects on an interface channel.
type Receiver struct {
	In  <-chan []byte
	Out chan<- interface{}
	serial.TypeDeserializer
	serial.TypeRegistrar
	errHandler lerr.ErrHandler
}

// NewListener creates a Listener that reads from the in channel,
// deserializes to a value and passes the value to a ListenerMuxer.
func NewListener(in <-chan []byte, d serial.TypeDeserializer, r serial.TypeRegistrar, errHandler any, handlers ...interface{}) (bus.Listener, error) {
	rc := &Receiver{
		In:               in,
		TypeDeserializer: d,
		TypeRegistrar:    r,
	}
	// TODO: don't use arbitrary size
	return listener.New(10, rc, errHandler, handlers...)
}

// Run starts the Receiver. It must be running to receive messages.
func (r *Receiver) Run() {
	for b := range r.In {
		r.handle(b)
	}
}

func (r *Receiver) handle(b []byte) {
	i, err := r.DeserializeType(b)
	if err != nil {
		r.errHandler.Handle(err)
		return
	}
	r.Out <- i
}

// SetOut sets the outgoing interface channel.
func (r *Receiver) SetOut(out chan<- interface{}) {
	r.Out = out
}

// SetErrorHandler to errHandler if ErrHandler is currently nil.
func (r *Receiver) SetErrorHandler(errHandler any) (err error) {
	r.errHandler, err = lerr.HandlerFunc(errHandler)
	return
}

func (r *Receiver) RegisterType(zeroValue interface{}) error {
	if r.TypeRegistrar == nil {
		return nil
	}
	return r.TypeRegistrar.RegisterType(zeroValue)
}
