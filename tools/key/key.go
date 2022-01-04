// Package key generates random keys, such as the ones used to sign sessions, and
// encodes them as text or as Go code.
package key

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/adamcolton/luce/lerr"
)

// reader fills a slice with random bytes. It is crypto/rand.Read, and it is a
// variable so that a test can make it fail.
var reader = rand.Read

// Key wraps a slice of bytes to help with generating and encoding.
type Key []byte

// DefaultLength is the length that will be generated if no length is given to
// New.
const DefaultLength = 32

// New creates a key using crypto/rand of the specified length. If the length
// is 0 then the default length is used. It panics if crypto/rand fails, because
// a key that was not filled with random bytes is not safe to use.
func New(ln int) Key {
	if ln <= 0 {
		ln = DefaultLength
	}
	key := make([]byte, ln)
	_, err := reader(key)
	lerr.Panic(err)
	return Key(key)
}

// Parse creates a key from a base64 URLEncoded string.
func Parse(s string) (Key, error) {
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
