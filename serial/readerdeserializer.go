package serial

import (
	"bytes"
	"io"
)

// ReaderDeserializer serializes the provided interface to the Reader.
type ReaderDeserializer func(interface{}, io.Reader) error

// Deserialize the interface from the byte slice.
func (fn ReaderDeserializer) Deserialize(i interface{}, b []byte) error {
	buf := bytes.NewBuffer(b)
	return fn(i, buf)
}
