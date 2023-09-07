// Package gob adapts encoding/gob to the serial interfaces.
package gob

import (
	"encoding/gob"
	"io"
)

// Register wraps gob.Register, so a caller that already imports this package does
// not also need to import encoding/gob to register a type.
func Register(value any) {
	gob.Register(value)
}

// Serialize encodes i with gob and writes it to w. It fulfills
// serial.WriterSerializer. A new Encoder is used for every call, so the output is
// self-contained and includes gob's type information.
func Serialize(i any, w io.Writer) error {
	return gob.NewEncoder(w).Encode(i)
}

// Deserialize reads gob data from r and decodes it into i. It fulfills
// serial.ReaderDeserializer.
func Deserialize(i any, r io.Reader) error {
	return gob.NewDecoder(r).Decode(i)
}
