package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	k := New(0)
	assert.Len(t, k, DefaultLength)
}

func TestCode(t *testing.T) {
	restore := reader
	reader = func(b []byte) (int, error) {
		for i := range b {
			b[i] = byte(i)
		}
		return len(b), nil
	}
	defer func() {
		reader = restore
	}()

	assert.Equal(t, "[]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}", New(0).Code())
}
