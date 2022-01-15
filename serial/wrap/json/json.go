package json

import (
	"bytes"
	"encoding/json"
	"io"
)

// Serialize wraps encoding/json Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(i interface{}, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/json Decoder to fulfill
// type32.DeserializeTypeID32Funcn
func Deserialize(i interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

type Serializer struct {
	Prefix, Indent string
}

func NewSerializer(prefix, indent string) Serializer {
	return Serializer{
		Prefix: prefix,
		Indent: indent,
	}
}

func (s Serializer) Serialize(i interface{}, buf []byte) ([]byte, error) {
	b := bytes.NewBuffer(buf)
	enc := json.NewEncoder(b)
	enc.SetIndent(s.Prefix, s.Indent)
	err := enc.Encode(i)
	return b.Bytes(), err
}

type Deserializer struct{}

func (Deserializer) Deserialize(v interface{}, data []byte) error {
	return json.Unmarshal(data, v)
}
