package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"

	"github.com/testify/assert"
)

func TestExternalType(t *testing.T) {
	ref := MustExternalPackageRef("foo")
	bar := ref.MustExternalType("Bar")

	str := PrefixWriteToString(bar, DefaultPrefixer)
	assert.Equal(t, "foo.Bar", str)
	assert.Equal(t, TypeDefKind, bar.Kind())
	assert.Equal(t, ref, bar.ExternalPackageRef())
	assert.Equal(t, PackageRef(ref), bar.PackageRef())

	i := NewImports(nil)
	bar.RegisterImports(i)
	str = bufpool.MustWriterToString(i)
	assert.Contains(t, str, "foo")

	_, err := ref.ExternalType("bar")
	assert.Equal(t, `ExternalType "bar" in package "foo" is not exported`, err.Error())
}
