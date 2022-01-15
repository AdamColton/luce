package json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

// Serialize wraps encoding/json Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(i any, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/json Decoder to fulfill
// type32.DeserializeTypeID32Funcn
func Deserialize(i any, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

// Serializer holds options for serializing json.
type Serializer struct {
	Prefix, Indent string
}

// NewSerializer with options
func NewSerializer(prefix, indent string) Serializer {
	return Serializer{
		Prefix: prefix,
		Indent: indent,
	}
}

// Serialize writes the JSON value of v to a byte slice. It will append to b.
func (s Serializer) Serialize(v any, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := s.WriteTo(v, buf)
	return buf.Bytes(), err
}

// WriteTo fulfills io.WriterTo. It writes the JSON value of v to w.
func (s Serializer) WriteTo(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent(s.Prefix, s.Indent)
	return enc.Encode(v)
}

var osCreate = func(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Save writes the JSON value of v to the file at the path location.
func (s Serializer) Save(v any, path string) error {
	f, err := osCreate(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.WriteTo(v, f)
}

// Deserializer fulfill serial.Deserializer.
type Deserializer struct{}

// Deserialize the JSON data into v. This fulfills fulfill serial.Deserializer.
func (Deserializer) Deserialize(v any, data []byte) error {
	return json.Unmarshal(data, v)
}

// ReadFrom read JSON data from r and deserializes it into v.
func (Deserializer) ReadFrom(v any, r io.Reader) error {
	return json.NewDecoder(r).Decode(v)
}

var osOpen = func(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// Load reads the JSON data from the file at the path location and deserializes
// it to v.
func (d Deserializer) Load(v any, path string) error {
	f, err := osOpen(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return d.ReadFrom(v, f)
}
