package serialbus

import (
	"errors"

	"github.com/adamcolton/luce/ds/bus"

	"github.com/adamcolton/luce/serial"
)

// Sender combines the logic of serializing an object and placing it
// on a channel
type Sender struct {
	Chan chan<- []byte
	serial.TypeSerializer
}

// Send takes a message, serializes it and places it on a channel.
func (s *Sender) Send(msg interface{}) error {
	b, err := s.SerializeType(msg, nil)
	if err != nil {
		return err
	}
	s.Chan <- b
	return nil
}

// MultiSender allows one message to be sent to multiple channels.
type MultiSender struct {
	Chans map[string]chan<- []byte
	serial.TypeSerializer
}

func NewMultiSender(s serial.TypeSerializer) *MultiSender {
	return &MultiSender{
		Chans:          make(map[string]chan<- []byte),
		TypeSerializer: s,
	}
}

// Send a message to the keys provided. If no keys are provided, the message will
// be sent to all channels.
func (ms *MultiSender) Send(msg interface{}, keys ...string) error {
	b, err := ms.SerializeType(msg, nil)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		for _, ch := range ms.Chans {
			ch <- b
		}
	} else {
		for _, key := range keys {
			if ch, found := ms.Chans[key]; found {
				ch <- b
			}
		}
	}

	return nil
}

// Add a chan<- []byte to the MultiSender and associate it with the key.
// If to is not of type chan<- []byte, an error is returned.
func (ms *MultiSender) Add(key string, to interface{}) error {
	ch, ok := to.(chan<- []byte)
	if !ok {
		bch, ok := to.(chan []byte)
		if !ok {
			return errors.New("Expected chan<- []byte")
		}
		ch = bch
	}
	ms.Chans[key] = ch
	return nil
}

// AddCh adds a chan<- []byte to the MultiSender and associate it with the
// key.
func (ms *MultiSender) AddCh(key string, ch chan<- []byte) {
	ms.Chans[key] = ch
}

// Delete a channel by key from the MultiSender.
func (ms *MultiSender) Delete(key string) {
	delete(ms.Chans, key)
}

// Receiver takes serialized messages off a serial bus, deserializes them
// and places the deserialized objects on an interface channel.
type Receiver struct {
	In  <-chan []byte
	Out chan<- interface{}
	serial.TypeDeserializer
	serial.TypeRegistrar
	ErrHandler func(error)
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
		if r.ErrHandler != nil {
			r.ErrHandler(err)
		}
		return
	}
	r.Out <- i
}

// SetOut sets the outgoing interface channel.
func (r *Receiver) SetOut(out chan<- interface{}) {
	r.Out = out
}

// SetErrorHandler to errHandler if ErrHandler is currently nil.
func (r *Receiver) SetErrorHandler(errHandler func(error)) {
	if r.ErrHandler == nil {
		r.ErrHandler = errHandler
	}
}

func (r *Receiver) RegisterType(zeroValue interface{}) error {
	if r.TypeRegistrar == nil {
		return nil
	}
	return r.TypeRegistrar.RegisterType(zeroValue)
}

// NewListener creates a Listener that reads from the in channel,
// deserializes to a value and passes the value to a ListenerMuxer.
func NewListener(in <-chan []byte, d serial.TypeDeserializer, r serial.TypeRegistrar, errHandler func(error), handlers ...interface{}) (bus.Listener, error) {
	rc := &Receiver{
		In:               in,
		TypeDeserializer: d,
		TypeRegistrar:    r,
	}
	lm, err := bus.NewListenerMux(nil, errHandler)
	if err != nil {
		return nil, err
	}
	return bus.NewListener(rc, lm, errHandler, handlers...)
}

// String converts []byte to string on a channel.
func String(in <-chan []byte) <-chan string {
	out := make(chan string, len(in))
	go func() {
		for b := range in {
			out <- string(b)
		}
		close(out)
	}()
	return out
}
