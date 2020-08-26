package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestArray(t *testing.T) {
	arr := IntType.Array(5)

	assert.Equal(t, IntType, arr.Elem())
	assert.Equal(t, 5, arr.Size)

	arr = IntType.Array(-5)
	str := PrefixWriteToString(arr, DefaultPrefixer)
	assert.Equal(t, "[...]int", str)
	assert.Equal(t, 0, arr.Size)
}
