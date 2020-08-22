package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestNewPackage(t *testing.T) {
	ctx := NewMemoryContext()
	pkg, err := ctx.Package("foo")
	assert.NoError(t, err)

	pkg2 := ctx.MustPackage("foo")
	assert.Equal(t, pkg, pkg2)

	assert.Equal(t, ErrBadImportPath, pkg.SetImportPath("bad import path"))
	assert.NoError(t, pkg.SetImportPath("bar"))
	assert.Equal(t, `"bar/foo"`, pkg.ImportSpec())

	ctx.Export()
}

func TestPkgErrors(t *testing.T) {
	_, err := NewPackage(nil, "foo")
	assert.Equal(t, ErrNilContext, err)

	ctx := NewMemoryContext()
	_, err = ctx.Package("bad package name")
	assert.Equal(t, ErrBadPackageName, err)

}
