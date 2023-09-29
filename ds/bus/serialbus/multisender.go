package serialbus

import (
	"errors"

	"github.com/adamcolton/luce/serial"
)

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
