package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/util/luceio"

	"github.com/testify/assert"
)

func TestTypeDef(t *testing.T) {
	ctx := NewMemoryContext()
	foo := ctx.MustPackage("baz").MustTypeDef("Foo", StringType)
	foo.Comment = "testing TypeDef comment"

	m := foo.MustMethod("Bar", IntType.Named("a")).
		Returns(BoolType.Named("isOne")).
		BodyString("isOne = a == 1\nreturn")
	m.Comment = "testing method comment"

	bar, ok := foo.Method("Bar")
	assert.True(t, ok)
	assert.Equal(t, m, bar)
	assert.NoError(t, m.Rename("Hide"))
	bar, ok = foo.Method("Bar")
	assert.False(t, ok)
	assert.Nil(t, bar)
	bar, ok = foo.Method("Hide")
	assert.True(t, ok)
	assert.Equal(t, m, bar)
	assert.NoError(t, m.Rename("Bar"))

	k := foo.File().MustTypeDef("Klaatu", IntType)
	k.Ptr = true
	str := PrefixWriteToString(k, DefaultPrefixer)
	assert.Equal(t, "*baz.Klaatu", str)

	foo.MustMethod("SingleReturn").
		UnnamedRets(BoolType).
		BodyString(`return false`).
		Ptr = true
	foo.MustMethod("MultiReturn").
		UnnamedRets(StringType, BoolType).
		BodyWriterTo(luceio.StringWriterTo(`return "hi", false`))

	ctx.MustExport()
	assert.Contains(t, ctx.Last(), "// Foo testing TypeDef comment\ntype Foo string")
	assert.Contains(t, ctx.Last(), "// Bar testing method comment\nfunc (f Foo) Bar(a int) (isOne bool) {\n\tisOne = a == 1\n\treturn\n}")
	assert.Contains(t, ctx.Last(), "func (f Foo) MultiReturn() (string, bool) {\n\treturn \"hi\", false\n}")
	assert.Contains(t, ctx.Last(), "func (f *Foo) SingleReturn() bool {\n\treturn false\n}")
}

func TestTypeDefRegisterImports(t *testing.T) {
	ctx := NewMemoryContext()
	foo := ctx.MustTypeDef("Foo", MustExternalPackageRef("bar").MustExternalType("Bar"))
	m := foo.MustMethod("ImportTest", MustExternalPackageRef("baz").MustExternalType("Baz").Named("b"))
	m.Body = testFuncBodyWriter("importTest.SayHi()")

	ctx.MustExport()
	assert.Contains(t, ctx.Last(), `"bar"`)
	assert.Contains(t, ctx.Last(), `"baz"`)
	assert.Contains(t, ctx.Last(), `"importTest"`)

}
