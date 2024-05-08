package gob

import (
	"bytes"
	"encoding/gob"
	"io"
)

// Serialize wraps encoding/gob Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(i interface{}, w io.Writer) error {
	return gob.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/gob Decoder to fulfill
// type32.DeserializeTypeID32Func
func Deserialize(i interface{}, r io.Reader) error {
	return gob.NewDecoder(r).Decode(i)
}

type Serializer struct{}

func (Serializer) Serialize(v any, buf []byte) ([]byte, error) {
	b := bytes.NewBuffer(buf[:0])
	err := Serialize(v, b)
	return b.Bytes(), err
}

type Deserializer struct{}

func (Deserializer) Deserialize(v any, data []byte) error {
	return Deserialize(v, bytes.NewBuffer(data))
}
