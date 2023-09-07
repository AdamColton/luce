package gob

import (
	"encoding/gob"
	"io"
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
