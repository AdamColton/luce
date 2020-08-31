package gothicgo

import (
	"io"
	"testing"

	"github.com/adamcolton/luce/util/luceio"

	"github.com/testify/assert"
)

type testFuncBodyWriter string

func (fw testFuncBodyWriter) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	i, err := w.Write([]byte(fw))
	return int64(i), err
}

func (fw testFuncBodyWriter) RegisterImports(i *Imports) {
	i.Add(MustPackageRef("importTest"))
}

func TestFunc(t *testing.T) {
	ctx, file := newFile("foo")
	args := []NameType{
		IntType.Named("a"),
		StringType.Named("b"),
		MustPackageRef("baz").NewTypeRef("Baz", nil).Named("c"),
	}
	fn := file.MustFunc("Rename", args...).
		UnnamedRets(BoolType)
	fn.Body = testFuncBodyWriter("return true")
	fn.Comment = "is a test function"
	assert.Equal(t, file, fn.File())
	fn.Rename("Bar")

	file.MustFunc("BodyStringTest", args...).
		BodyString("return bodystring").
		UnnamedRets(BoolType)

	file.MustFunc("BodyWriterToTest", args...).
		BodyWriterTo(luceio.StringWriterTo("return bodywriterto")).
		UnnamedRets(BoolType)

	ctx.MustExport()

	assert.Equal(t, "func Bar(a int, b string, c baz.Baz) bool", PrefixWriteToString(fn.Ref(), DefaultPrefixer))

	assert.Contains(t, ctx.Last(), "func Bar(a int, b string, c baz.Baz) bool {\n\treturn true\n}")
	assert.Contains(t, ctx.Last(), "func BodyStringTest(a int, b string, c baz.Baz) bool {\n\treturn bodystring\n}")
	assert.Contains(t, ctx.Last(), "func BodyWriterToTest(a int, b string, c baz.Baz) bool {\n\treturn bodywriterto\n}")
	assert.Contains(t, ctx.Last(), "importTest")
	assert.Contains(t, ctx.Last(), "// Bar is a test function")

	assert.Equal(t, "foo.Bar(x, y, z)", fn.Call(DefaultPrefixer, "x", "y", "z"))
}

func TestNewFuncErr(t *testing.T) {
	_, file := newFile("foo")
	args := []NameType{
		IntType.Named("a"),
		StringType.Unnamed(),
		MustPackageRef("baz").NewTypeRef("Baz", nil).Named("c"),
	}
	_, err := file.NewFunc("someFunc", args...)
	assert.Equal(t, ErrUnnamedFuncArg, err)
}
