package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestPointer(t *testing.T) {
	ptr := IntType.Ptr()
	str := PrefixWriteToString(ptr, DefaultPrefixer)

	assert.Equal(t, "*int", str)
	assert.Equal(t, PointerKind, ptr.Kind())
	assert.Equal(t, IntType, ptr.Elem())
	assert.Equal(t, IntType, ptr.PointerElem())
	assert.Equal(t, PkgBuiltin(), ptr.PackageRef())
}
