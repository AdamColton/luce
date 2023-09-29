package serialbus

import "github.com/adamcolton/luce/serial"

// Sender combines the logic of serializing an object and placing it
// on a channel. It fulfills bus.Sender.
type Sender struct {
	// Chan is the serial bus that the messages are sent on.
	Chan chan<- []byte
	serial.TypeSerializer
}

// Send takes a message, serializes it with its type and places it on the channel.
// It waits until the channel can take it. If the message can't be serialized,
// it returns the error and sends nothing.
func (s *Sender) Send(msg interface{}) error {
	b, err := s.SerializeType(msg, nil)
	if err != nil {
		return err
	}
	s.Chan <- b
	return nil
}
