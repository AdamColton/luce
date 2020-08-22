package gothicgo

import (
	"bytes"
	"testing"

	"github.com/testify/assert"
)

func TestPointer(t *testing.T) {
	ptr := IntType.Ptr()
	buf := bytes.NewBuffer(nil)
	ptr.PrefixWriteTo(buf, DefaultPrefixer)

	assert.Equal(t, "*int", buf.String())
	assert.Equal(t, PointerKind, ptr.Kind())
	assert.Equal(t, IntType, ptr.Elem())
	assert.Equal(t, IntType, ptr.PointerElem())
	assert.Equal(t, PkgBuiltin(), ptr.PackageRef())
}
