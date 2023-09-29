package serialbus

import "github.com/adamcolton/luce/serial"

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
