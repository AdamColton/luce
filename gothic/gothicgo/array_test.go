package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestArray(t *testing.T) {
	arr := IntType.Array(5)
	str := PrefixWriteToString(arr, DefaultPrefixer)

	assert.Equal(t, "[5]int", str)
	assert.Equal(t, ArrayKind, arr.Kind())
	assert.Equal(t, IntType, arr.Elem())
	assert.Equal(t, 5, arr.Size)
	assert.Equal(t, PkgBuiltin(), arr.PackageRef())

	n := arr.Named("Foo")
	assert.Equal(t, "Foo", n.Name())
	assert.Equal(t, arr, n.Type)

	n = arr.Unnamed()
	assert.Equal(t, "", n.Name())
	assert.Equal(t, arr, n.Type)

	p := arr.Ptr()
	assert.Equal(t, PointerKind, p.Kind())
	assert.Equal(t, arr, p.Elem())

	s := arr.Slice()
	assert.Equal(t, SliceKind, s.Kind())
	assert.Equal(t, arr, s.Elem())

	a := arr.Array(13)
	assert.Equal(t, ArrayKind, a.Kind())
	assert.Equal(t, arr, a.Elem())

	mp := arr.AsMapElem(IntType)
	assert.Equal(t, MapKind, mp.Kind())
	assert.Equal(t, arr, mp.Elem())

	mp = arr.AsMapKey(IntType)
	assert.Equal(t, MapKind, mp.Kind())
	assert.Equal(t, arr, mp.MapKey())

	arr = IntType.Array(-5)
	str = PrefixWriteToString(arr, DefaultPrefixer)
	assert.Equal(t, "[...]int", str)
	assert.Equal(t, 0, arr.Size)
}
