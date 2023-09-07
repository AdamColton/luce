package json

import (
	"encoding/json"
	"io"
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
