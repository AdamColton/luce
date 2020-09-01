package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/testify/assert"
)

func TestPackageVar(t *testing.T) {
	ctx, file := newFile("foo")

	_, err := file.NewPackageVar("Bar", MustPackageRef("bar").NewTypeRef("Bar", nil))
	assert.NoError(t, err)

	pv, err := file.NewPackageVar("A", StringType)
	pv.Value = IgnorePrefixer{luceio.StringWriterTo(`"this is a test"`)}
	assert.NoError(t, err)

	pv, err = file.NewPackageVar("B", nil)
	pv.Value = MustPackageRef("importtest").
		NewFuncRef("FuncCall", StringType.Unnamed()).
		Caller(`"test"`)

	assert.NoError(t, err)

	assert.Equal(t, file.pkg, pv.PackageRef())
	assert.Equal(t, file, pv.File())

	_, err = file.NewPackageVar("B", nil)
	assert.Equal(t, "NewPackageVar: File.AddGenerator: Name 'B' already exists in package 'foo'", err.Error())

	ctx.MustExport()

	assert.Contains(t, ctx.Last(), "var Bar bar.Bar")
	assert.Contains(t, ctx.Last(), `"bar"`)
	assert.Contains(t, ctx.Last(), `"importtest"`)
	assert.Contains(t, ctx.Last(), "var A string = \"this is a test\"")
	assert.Contains(t, ctx.Last(), `var B = importtest.FuncCall("test")`)

	r := pv.Ref()
	str := PrefixWriteToString(r, DefaultPrefixer)
	assert.Equal(t, "foo.B", str)

	i := NewImports(nil)
	r.RegisterImports(i)
	str = bufpool.MustWriterToString(i)
	assert.Contains(t, str, `"foo"`)
	assert.Equal(t, file.Package(), r.PackageRef())
}
