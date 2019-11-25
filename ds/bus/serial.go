package bus

import (
	"errors"

	"github.com/adamcolton/luce/serial"
)

// SerialSender combines the logic of serializing an object and placing it
// on a channel
type SerialSender struct {
	Ch chan<- []byte
	serial.Serializer
}

// Send takes a message, serializes it and places it on a channel.
func (s *SerialSender) Send(msg interface{}) error {
	b, err := s.Serialize(msg)
	if err != nil {
		return err
	}
	s.Ch <- b
	return nil
}

// SerialMultiSender allows one message to be sent to multiple channels.
type SerialMultiSender struct {
	Chans map[string]chan<- []byte
	serial.Serializer
}

func NewSerialMultiSender(s serial.Serializer) *SerialMultiSender {
	return &SerialMultiSender{
		Chans:      make(map[string]chan<- []byte),
		Serializer: s,
	}
}

// Send a message to the keys provided. If no keys are provided, the message will
// be sent to all channels.
func (ms *SerialMultiSender) Send(msg interface{}, keys ...string) error {
	b, err := ms.Serialize(msg)
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

// Add a chan<- []byte to the SerialMultiSender and associate it with the key.
// If to is not of type chan<- []byte, an error is returned.
func (ms *SerialMultiSender) Add(key string, to interface{}) error {
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

// AddCh adds a chan<- []byte to the SerialMultiSender and associate it with the
// key.
func (ms *SerialMultiSender) AddCh(key string, ch chan<- []byte) {
	ms.Chans[key] = ch
}

// Delete a channel by key from the SerialMultiSender.
func (ms *SerialMultiSender) Delete(key string) {
	delete(ms.Chans, key)
}

// SerialReceiver takes serialized messages off a serial bus, deserializes them
// and places the deserialized objects on an interface channel.
type SerialReceiver struct {
	In  <-chan []byte
	Out chan<- interface{}
	serial.Deserializer
	ErrHandler func(error)
}

// Run starts the SerialReceiver. It must be running to receive messages.
func (r *SerialReceiver) Run() {
	for b := range r.In {
		r.handle(b)
	}
}

func (r *SerialReceiver) handle(b []byte) {
	i, err := r.Deserialize(b)
	if err != nil {
		if r.ErrHandler != nil {
			r.ErrHandler(err)
		}
		return
	}
	r.Out <- i
}

// SetOut sets the outgoing interface channel.
func (r *SerialReceiver) SetOut(out chan<- interface{}) {
	r.Out = out
}

// SetErrorHandler to errHandler if ErrHandler is currently nil.
func (r *SerialReceiver) SetErrorHandler(errHandler func(error)) {
	if r.ErrHandler == nil {
		r.ErrHandler = errHandler
	}
}

// NewSerialListener creates a Listener that reads from the in channel,
// deserializes to a value and passes the value to a ListenerMuxer.
func NewSerialListener(in <-chan []byte, d serial.Deserializer, errHandler func(error), handlers ...interface{}) (Listener, error) {
	r := &SerialReceiver{
		In:           in,
		Deserializer: d,
	}
	lm, err := NewListenerMux(nil, errHandler)
	if err != nil {
		return nil, err
	}
	return NewListener(r, lm, errHandler, handlers...)
}
