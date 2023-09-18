package serial

import (
	"bytes"
	"io"
)

// ReaderDeserializer serializes the provided interface to the Reader.
type ReaderDeserializer func(any, io.Reader) error

// Deserialize the interface from the byte slice.
func (fn ReaderDeserializer) Deserialize(i any, b []byte) error {
	buf := bytes.NewBuffer(b)
	return fn(i, buf)
}
