package serial

import (
	"bytes"
	"io"
)

// ReaderDeserializer is a function that populates a value from an io.Reader. It
// fulfills Deserializer.
type ReaderDeserializer func(any, io.Reader) error

// Deserialize calls fn with a reader over b.
func (fn ReaderDeserializer) Deserialize(i any, b []byte) error {
	buf := bytes.NewBuffer(b)
	return fn(i, buf)
}
