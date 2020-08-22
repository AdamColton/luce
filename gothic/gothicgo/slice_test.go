package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestSlice(t *testing.T) {
	slc := IntType.Slice()
	str := PrefixWriteToString(slc, DefaultPrefixer)

	assert.Equal(t, "[]int", str)
	assert.Equal(t, SliceKind, slc.Kind())
	assert.Equal(t, IntType, slc.Elem())
	assert.Equal(t, IntType, slc.SliceElem())
	assert.Equal(t, PkgBuiltin(), slc.PackageRef())
}
