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

func TestNewHelpfulTypeWrapper(t *testing.T) {
	doubleWrap := HelpfulTypeWrapper{IntType}

	// confirm doubleWrap
	_, ok := doubleWrap.Type.(HelpfulTypeWrapper)
	assert.True(t, ok)

	singleWrap := NewHelpfulTypeWrapper(doubleWrap)

	// confirm singleWrap
	_, ok = singleWrap.Type.(HelpfulTypeWrapper)
	assert.False(t, ok)
}

func TestNewHelpfulType(t *testing.T) {
	singleWrap := NewHelpfulType(IntType)

	// confirm singleWrap was not re-wrapped
	swht, ok := singleWrap.(HelpfulTypeWrapper)
	assert.True(t, ok)
	_, ok = swht.Type.(HelpfulTypeWrapper)
	assert.False(t, ok)
}
