package json32

import (
	"encoding/json"
	"io"
)

// Serialize wraps encoding/json Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(w io.Writer, i interface{}) error {
	return json.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/json Decoder to fulfill
// type32.DeserializeTypeID32Func
func Deserialize(r io.Reader, i interface{}) error {
	return json.NewDecoder(r).Decode(i)
}
