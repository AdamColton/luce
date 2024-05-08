// Package gob adapts encoding/gob to the serial interfaces.
package gob

import (
	"bytes"
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

// Encoder creates an instance of gob.Encoder using a bytes.Buffer constructed
// with b. Both the encoder and the buffer are returned.
func Encoder(b []byte) (*gob.Encoder, *bytes.Buffer) {
	buf := bytes.NewBuffer(b)
	return gob.NewEncoder(buf), buf
}

// Decoder creates an instance of gob.Decoder using a bytes.Buffer constructed
// with data.
func Decoder(data []byte) *gob.Decoder {
	return gob.NewDecoder(bytes.NewBuffer(data))
}

// Serializer fulfills serial.Serializer with gob.
type Serializer struct{}

// Serialize v using gob and append it to b. Every call uses a new Encoder, so the
// result is self-contained.
func (Serializer) Serialize(v any, b []byte) ([]byte, error) {
	enc, buf := Encoder(b)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

// Deserializer fulfills serial.Deserializer with gob.
type Deserializer struct{}

// Deserialize decodes data into v using gob.
func (Deserializer) Deserialize(v any, data []byte) error {
	return Decoder(data).Decode(v)
}
