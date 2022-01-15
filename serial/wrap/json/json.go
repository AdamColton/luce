// Package json adapts encoding/json to the serial interfaces.
package json

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

// Serialize encodes i as JSON and writes it to w. It fulfills
// serial.WriterSerializer.
func Serialize(i any, w io.Writer) error {
	return json.NewEncoder(w).Encode(i)
}

// Deserialize reads JSON from r and decodes it into i. It fulfills
// serial.ReaderDeserializer.
func Deserialize(i any, r io.Reader) error {
	return json.NewDecoder(r).Decode(i)
}

// Serializer serializes values as JSON. Prefix and Indent are passed to
// json.Encoder.SetIndent, so when both are empty the output is not indented.
type Serializer struct {
	Prefix, Indent string
}

// NewSerializer creates a Serializer that indents with the given prefix and
// indent.
func NewSerializer(prefix, indent string) Serializer {
	return Serializer{
		Prefix: prefix,
		Indent: indent,
	}
}

// Serialize writes the JSON value of v to a byte slice. It will append to b. It
// fulfills serial.Serializer.
func (s Serializer) Serialize(v any, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := s.WriteTo(v, buf)
	return buf.Bytes(), err
}

// WriteTo writes the JSON value of v to w. It fulfills serial.WriterSerializer.
func (s Serializer) WriteTo(v any, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent(s.Prefix, s.Indent)
	return enc.Encode(v)
}

var osCreate = func(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Save writes the JSON value of v to the file at the path location. The file is
// created, or truncated if it already exists.
func (s Serializer) Save(v any, path string) error {
	f, err := osCreate(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.WriteTo(v, f)
}

// Deserializer deserializes JSON. It fulfills serial.Deserializer.
type Deserializer struct{}

// Deserialize unmarshals the JSON data into v.
func (Deserializer) Deserialize(v any, data []byte) error {
	return json.Unmarshal(data, v)
}

// ReadFrom reads JSON data from r and deserializes it into v. It fulfills
// serial.ReaderDeserializer.
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
