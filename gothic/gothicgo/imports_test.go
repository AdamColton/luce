package gothicgo

import (
	"strings"
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"
	"github.com/stretchr/testify/assert"
)

var toStr = bufpool.MustWriterToString

func TestImportAdd(t *testing.T) {
	i := NewImports(nil)
	assert.Equal(t, "", toStr(i))

	i.Add(MustPackageRef("foo/bar"), nil, MustPackageRef("foo/baz"))
	assert.Equal(t, "import (\n\t\"foo/bar\"\n\t\"foo/baz\"\n)\n", toStr(i))
}

func TestImportDouble(t *testing.T) {
	i := NewImports(nil)
	pkg1 := MustPackageRef("foo/bar")
	pkg2 := MustPackageRef("foo/bar")
	i.Add(pkg1)
	i.Add(pkg2)
	assert.Equal(t, 1, strings.Count(toStr(i), "foo/bar"))
}

func TestImportAddImports(t *testing.T) {
	i1 := NewImports(nil)
	i1.Add(MustPackageRef("foo/bar"))

	i2 := NewImports(nil)
	i2.AddImports(i1)

	assert.Equal(t, "import (\n\t\"foo/bar\"\n)\n", toStr(i2))
}

func TestImportPrefix(t *testing.T) {
	var i *Imports
	bar := MustPackageRef("foo/bar")
	// i.Prefix should work even with nil *Imports
	assert.Equal(t, "bar.", i.Prefix(bar))
	assert.Equal(t, "", i.Prefix(PkgBuiltin()))

	i = NewImports(nil)
	i.Add(bar)

	assert.Equal(t, "bar.", i.Prefix(bar))

	baz := MustPackageRef("foo/baz")
	assert.Equal(t, "baz.", i.Prefix(baz))
}

func TestImportSelf(t *testing.T) {
	bar := MustPackageRef("foo/bar")
	i := NewImports(bar)
	i.Add(bar) //should be safe to add self
	assert.Equal(t, "", i.Prefix(bar))
}

func TestImportRemove(t *testing.T) {
	i := NewImports(nil)

	bar := MustPackageRef("foo/bar")
	baz := MustPackageRef("foo/baz")
	i.Add(bar, baz)
	i.RemoveRef(baz)
	assert.Equal(t, "import (\n\t\"foo/bar\"\n)\n", bufpool.MustWriterToString(i))
}

func TestAlias(t *testing.T) {
	i := NewImports(nil)
	bar := MustPackageRef("foo/bar")
	a := NewAlias(bar, "fb")
	i.Add(a)
	i.Add(bar) // should have no effect

	assert.Equal(t, "fb.", i.Prefix(bar))

	assert.Equal(t, "import (\n\tfb \"foo/bar\"\n)\n", bufpool.MustWriterToString(i))
}
