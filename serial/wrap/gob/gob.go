package gob

import (
	"encoding/gob"
	"io"
)

// Serialize wraps encoding/gob Encoder to fulfill type32.SerializeTypeID32Func
func Serialize(w io.Writer, i interface{}) error {
	return gob.NewEncoder(w).Encode(i)
}

// Deserialize wraps encoding/gob Decoder to fulfill
// type32.DeserializeTypeID32Func
func Deserialize(r io.Reader, i interface{}) error {
	return gob.NewDecoder(r).Decode(i)
}
