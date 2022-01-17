package json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
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
	return b.Bytes(), s.WriteTo(i, b)
}

func (s Serializer) WriteTo(i interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent(s.Prefix, s.Indent)
	return enc.Encode(i)
}

func (s Serializer) Save(i interface{}, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.WriteTo(i, f)
}

type Deserializer struct{}

func (Deserializer) Deserialize(v interface{}, data []byte) error {
	return json.Unmarshal(data, v)
}

func (Deserializer) ReadFrom(v interface{}, r io.Reader) error {
	return json.NewDecoder(r).Decode(v)
}

func (d Deserializer) Load(v interface{}, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return d.ReadFrom(v, f)
}
