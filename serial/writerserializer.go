package serial

import (
	"bytes"
	"io"
)

// WriterSerializer serializes the provided interface to the Writer.
type WriterSerializer func(interface{}, io.Writer) error

// Serialize the interface to the byte slice.
func (fn WriterSerializer) Serialize(i interface{}, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := fn(i, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
