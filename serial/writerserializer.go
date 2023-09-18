package serial

import (
	"bytes"
	"io"
)

// WriterSerializer is a function that serializes a value to an io.Writer. It
// fulfills Serializer.
type WriterSerializer func(any, io.Writer) error

// Serialize calls fn with a buffer that starts with b and returns the buffer's
// bytes, so the serialization is appended to b. On error it returns nil, not
// b.
func (fn WriterSerializer) Serialize(i any, b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	err := fn(i, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
