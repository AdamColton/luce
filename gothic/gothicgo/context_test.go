package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := NewMemoryContext()

	assert.NoError(t, ctx.Prepare())
	assert.NoError(t, ctx.Generate())

	ctx.SetOutputPath("foo")
	assert.Equal(t, "foo/bar", ctx.OutputPath("bar"))

	assert.NoError(t, ctx.SetImportPath("baz"))
	assert.Equal(t, "baz", ctx.ImportPath())
	assert.Error(t, ctx.SetImportPath("bad import path"))

}
