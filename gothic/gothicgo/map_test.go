package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestMap(t *testing.T) {
	mp := IntType.AsMapElem(StringType)
	str := PrefixWriteToString(mp, DefaultPrefixer)

	assert.Equal(t, "map[string]int", str)
	assert.Equal(t, MapKind, mp.Kind())
	assert.Equal(t, PkgBuiltin(), mp.PackageRef())
	assert.Equal(t, IntType, mp.Elem())
	assert.Equal(t, IntType, mp.MapElem())
	assert.Equal(t, StringType, mp.MapKey())

	mp = IntType.AsMapKey(StringType)
	str = PrefixWriteToString(mp, DefaultPrefixer)
	assert.Equal(t, "map[int]string", str)
	assert.Equal(t, StringType, mp.Elem())
	assert.Equal(t, StringType, mp.MapElem())
	assert.Equal(t, IntType, mp.MapKey())
}

func TestMapRegisterImports(t *testing.T) {
	i := NewImports(nil)
	mp := IntType.AsMapKey(StringType)

	mp.RegisterImports(i)

	// todo: after external type defs are done use types that will cause
	// registration.
}
