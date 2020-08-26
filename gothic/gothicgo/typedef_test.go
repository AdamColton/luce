package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestTypeDef(t *testing.T) {
	ctx := NewMemoryContext()
	foo, err := ctx.NewTypeDef("Foo", StringType)
	assert.NoError(t, err)

	args := []NameType{
		IntType.Named("a"),
	}
	rets := []NameType{
		BoolType.Unnamed(),
	}
	m, err := foo.NewMethod("Bar", args, rets, false)
	assert.NoError(t, err)
	m.BodyString("return false")

	ctx.MustExport()
	assert.Contains(t, ctx.Last(), "type Foo string")
	assert.Contains(t, ctx.Last(), "func (f Foo) Bar(a int) bool {\n\treturn false\n}")
}
