// Package json adapts encoding/json to the serial interfaces.
package json

import (
	"encoding/json"
	"io"
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
