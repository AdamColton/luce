package gothicgo

import (
	"bytes"
	"testing"

	"github.com/testify/assert"
)

func TestBuiltinType(t *testing.T) {
	assert.Equal(t, BoolKind, BoolType.Kind())
	assert.Equal(t, PkgBuiltin(), BoolType.PackageRef())

	buf := bytes.NewBuffer(nil)
	BoolType.PrefixWriteTo(buf, DefaultPrefixer)
	assert.Equal(t, "bool", buf.String())

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
