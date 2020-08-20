package gothicgo

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportAdd(t *testing.T) {
	i := NewImports(nil)
	buf := bytes.NewBuffer(nil)
	i.WriteTo(buf)
	assert.Equal(t, "", buf.String())

	i.Add(MustPackageRef("foo/bar"), nil, MustPackageRef("foo/baz"))
	buf.Reset()
	i.WriteTo(buf)
	assert.Equal(t, "import (\n\t\"foo/bar\"\n\t\"foo/baz\"\n)\n", buf.String())
}

func TestImportDouble(t *testing.T) {
	i := NewImports(nil)
	pkg1 := MustPackageRef("foo/bar")
	pkg2 := MustPackageRef("foo/bar")
	i.Add(pkg1)
	i.Add(pkg2)
	assert.Equal(t, 1, strings.Count(i.String(), "foo/bar"))
}

func TestImportAddImports(t *testing.T) {
	i1 := NewImports(nil)
	i1.Add(MustPackageRef("foo/bar"))

	i2 := NewImports(nil)
	i2.AddImports(i1)

	assert.Equal(t, "import (\n\t\"foo/bar\"\n)\n", i2.String())
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
	buf := bytes.NewBuffer(nil)

	bar := MustPackageRef("foo/bar")
	baz := MustPackageRef("foo/baz")
	i.Add(bar, baz)
	i.RemoveRef(baz)
	i.WriteTo(buf)
	assert.Equal(t, "import (\n\t\"foo/bar\"\n)\n", buf.String())
}

func TestAlias(t *testing.T) {
	i := NewImports(nil)
	bar := MustPackageRef("foo/bar")
	a := NewAlias(bar, "fb")
	i.Add(a)
	i.Add(bar) // should have no effect

	assert.Equal(t, "fb.", i.Prefix(bar))

	buf := bytes.NewBuffer(nil)
	i.WriteTo(buf)
	assert.Equal(t, "import (\n\tfb \"foo/bar\"\n)\n", buf.String())
}
