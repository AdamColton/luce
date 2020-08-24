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
	i.Add(MustExternalPackageRef("importTest"))
}

func TestFunc(t *testing.T) {
	ctx, file := newFile("foo")
	args := []NameType{
		IntType.Named("a"),
		StringType.Named("b"),
		MustExternalPackageRef("baz").MustExternalType("Baz").Named("c"),
	}
	rets := []NameType{
		BoolType.Unnamed(),
	}
	fn := file.MustFunc("Rename", args, rets, false)
	fn.Body = testFuncBodyWriter("return true")
	fn.Comment = "is a test function"
	assert.Equal(t, file, fn.File())
	fn.Rename("Bar")

	file.MustFunc("BodyStringTest", args, rets, false).BodyString("return bodystring")
	file.MustFunc("BodyWriterToTest", args, rets, false).BodyWriterTo(luceio.StringWriterTo("return bodywriterto"))

	ctx.MustExport()

	str := ctx.Last.String()
	assert.Contains(t, str, "func Bar(a int, b string, c baz.Baz) bool {\n\treturn true\n}")
	assert.Contains(t, str, "func BodyStringTest(a int, b string, c baz.Baz) bool {\n\treturn bodystring\n}")
	assert.Contains(t, str, "func BodyWriterToTest(a int, b string, c baz.Baz) bool {\n\treturn bodywriterto\n}")
	assert.Contains(t, str, "importTest")
	assert.Contains(t, str, "// Bar is a test function")

	str = fn.Call(DefaultPrefixer, "x", "y", "z")
	assert.Equal(t, "foo.Bar(x, y, z)", str)
}
