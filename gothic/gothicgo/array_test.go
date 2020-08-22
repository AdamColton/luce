package gothicgo

import (
	"bytes"
	"testing"

	"github.com/testify/assert"
)

func TestArray(t *testing.T) {
	arr := IntType.Array(5)
	buf := bytes.NewBuffer(nil)
	arr.PrefixWriteTo(buf, DefaultPrefixer)

	assert.Equal(t, "[5]int", buf.String())
	assert.Equal(t, ArrayKind, arr.Kind())
	assert.Equal(t, IntType, arr.Elem())
	assert.Equal(t, IntType, arr.ArrayElem())
	assert.Equal(t, 5, arr.Size())

	arr = IntType.Array(-5)
	buf.Reset()
	arr.PrefixWriteTo(buf, DefaultPrefixer)
	assert.Equal(t, "[...]int", buf.String())
	assert.Equal(t, 0, arr.Size())
}
