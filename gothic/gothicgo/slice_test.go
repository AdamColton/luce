package gothicgo

import (
	"bytes"
	"testing"

	"github.com/testify/assert"
)

func TestSlice(t *testing.T) {
	slc := IntType.Slice()
	buf := bytes.NewBuffer(nil)
	slc.PrefixWriteTo(buf, DefaultPrefixer)

	assert.Equal(t, "[]int", buf.String())
	assert.Equal(t, SliceKind, slc.Kind())
	assert.Equal(t, IntType, slc.Elem())
	assert.Equal(t, IntType, slc.SliceElem())
}
