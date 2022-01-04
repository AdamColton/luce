package key

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"strconv"
)

var reader = rand.Read

// Key wraps a slice of bytes to help with generating and encoding.
type Key []byte

// DefaultLength is the length that will be generated if no length is given to
// New
const DefaultLength = 32

// New creates a key using crypto/rand of the specified length. If the length
// is 0 then the default length is used.
func New(ln int) Key {
	if ln <= 0 {
		ln = DefaultLength
	}
	key := make([]byte, ln)
	reader(key)
	return Key(key)
}

// String creates a key from a base64 URLEncoded string.
func String(s string) (Key, error) {
	return base64.URLEncoding.DecodeString(s)
}

// Code converts the key to a string of Go code.
func (k Key) Code() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[]byte{")
	for i, b := range k {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(strconv.Itoa(int(b)))
	}
	buf.WriteString("}")
	return buf.String()
}

// String fulfills stringer and encodes the key using base64 URLEncoding.
func (k Key) String() string {
	return base64.URLEncoding.EncodeToString(k)
}
