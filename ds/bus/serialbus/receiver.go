package serialbus

// == projects.Code.luce.serialbus ==

import (
	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/listener"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// Receiver takes serialized messages off a serial bus, deserializes them and
// places the deserialized values on an interface channel. It fulfills
// bus.Receiver.
type Receiver struct {
	// In is the serial bus that the messages are received from.
	In <-chan []byte
	// Out is the channel that the values are sent on. SetOut sets it.
	Out chan<- interface{}
	serial.TypeDeserializer
	serial.TypeRegistrar
	errHandler lerr.ErrHandler
}

// NewListener creates a Listener that reads serialized messages from the in
// channel, deserializes each to a value with d, and passes the value to the
// handler for its type. The types that the handlers take are registered with r.
// Errors, from deserializing and from the handlers' switch, go to errHandler,
// which may be nil. It has room for 10 handlers.
func NewListener(in <-chan []byte, d serial.TypeDeserializer, r serial.TypeRegistrar, errHandler any, handlers ...interface{}) (bus.Listener, error) {
	rc := &Receiver{
		In:               in,
		TypeDeserializer: d,
		TypeRegistrar:    r,
	}
	// [ ] serialbus.NewListener arbitrary size
	//	don't use arbitrary size
	return listener.New(10, rc, errHandler, handlers...)
}

// Run receives messages until In is closed, and sends each value to Out. It
// must be running to receive messages, and it waits while nothing is reading Out.
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

// SetErrorHandler sets what is called with each error from deserializing, either
// a func(error) or a channel of errors, as lerr.HandlerFunc accepts, or nil to
// ignore errors. If it returns an error the current handler is unchanged.
func (r *Receiver) SetErrorHandler(errHandler any) error {
	h, err := lerr.HandlerFunc(errHandler)
	if err != nil {
		return err
	}
	r.errHandler = h
	return nil
}

// RegisterType registers the type of zeroValue with the TypeRegistrar, so that
// messages of that type can be deserialized. It does nothing if there is no
// TypeRegistrar.
func (r *Receiver) RegisterType(zeroValue interface{}) error {
	if r.TypeRegistrar == nil {
		return nil
	}
	return r.TypeRegistrar.RegisterType(zeroValue)
}
