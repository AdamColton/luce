package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestInterface(t *testing.T) {
	fs := NewFuncSig("Foo")
	ir := MustPackageRef("bar").NewInterfaceRef("Bar")
	i := NewInterfaceType(fs, ir)

	str := PrefixWriteToString(i, DefaultPrefixer)
	assert.Equal(t, "interface {\n\tFoo()\n\tbar.Bar\n}", str)

	i = NewInterfaceType()
	str = PrefixWriteToString(i, DefaultPrefixer)
	assert.Equal(t, "interface{}", str)
}

func TestInterfaceDef(t *testing.T) {
	fs := NewFuncSig("Foo")
	ir := MustPackageRef("bar").NewInterfaceRef("Bar")

	ctx, file := newFile("baz")
	i, err := file.NewInterfaceDef("Baz", fs)
	assert.NoError(t, err)
	i.Embed(ir)

	ctx.MustExport()

	assert.Contains(t, ctx.Last(), "type Baz interface {\n\tFoo()\n\tbar.Bar\n}")
	assert.NotContains(t, ctx.Last(), `"foo"`)

	r := i.Ref()
	str := PrefixWriteToString(r, DefaultPrefixer)
	assert.Equal(t, "baz.Baz", str)
	assert.Equal(t, ctx.MustPackage("baz"), i.PackageRef())
	assert.Equal(t, file, i.File())
	assert.Equal(t, i.Interface, i.Elem())
}
