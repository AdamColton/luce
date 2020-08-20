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
}
