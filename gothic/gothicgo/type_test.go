package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestBuiltinType(t *testing.T) {
	assert.Equal(t, BoolKind, BoolType.Kind())
	assert.Equal(t, PkgBuiltin(), BoolType.PackageRef())

	str := PrefixWriteToString(BoolType, DefaultPrefixer)
	assert.Equal(t, "bool", str)

	assert.Equal(t, "foo bool", PrefixWriteToString(BoolType.Named("foo"), DefaultPrefixer))
	assert.Equal(t, "bool", PrefixWriteToString(BoolType.Unnamed(), DefaultPrefixer))
}
