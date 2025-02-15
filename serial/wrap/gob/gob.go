package gob

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/adamcolton/luce/lerr"
)

// Register wraps gob.Register to avoid importing this package and encoding/gob
// just to invoke register.
func Register(value any) {
	gob.Register(value)
}

// Serialize wraps encoding/gob Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(i any, w io.Writer) error {
	return gob.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/gob Decoder to fulfill
// type32.DeserializeTypeID32Func
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
// with b.
func Decoder(data []byte) *gob.Decoder {
	return gob.NewDecoder(bytes.NewBuffer(data))
}

// Serializer fulfills serial.Serializer with gob.
type Serializer struct{}

// Serialize v using gob and append it to b
func (Serializer) Serialize(v any, b []byte) ([]byte, error) {
	enc, buf := Encoder(b)
	err := enc.Encode(v)
	return buf.Bytes(), err
}

// Deserializer fulfills serial.Deserializer with gob.
type Deserializer struct{}

// Deserialize data to v using gob.
func (Deserializer) Deserialize(v any, data []byte) error {
	return Decoder(data).Decode(v)
}

// Enc encodes v to a []byte using gob. It will panic if there is an error.
func Enc(v any) []byte {
	enc, buf := Encoder(nil)
	lerr.Panic(enc.Encode(v))
	return buf.Bytes()
}

// Dec decodes data to using gob. It will panic if there is an error.
func Dec(data []byte, v any) {
	lerr.Panic(Decoder(data).Decode(v))
}
