package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestComment(t *testing.T) {
	ctx, file := newFile("foo")
	assert.Equal(t, defaultCommentWidth, file.CW)

	file.NewComment("This is a test")
	file.CW = 10
	file.NewComment("The sun was shining on the sea")

	assert.Equal(t, ErrCommentWidth, ctx.SetCommentWidth(0))

	ctx.MustExport()

	assert.Contains(t, ctx.Last(), "// This is a test")
	assert.Contains(t, ctx.Last(), "// The sun")
	assert.Contains(t, ctx.Last(), "// was")
	assert.Contains(t, ctx.Last(), "// shining")
	assert.Contains(t, ctx.Last(), "// on the")
	assert.Contains(t, ctx.Last(), "// sea")
}

func TestDefaultComment(t *testing.T) {
	ctx := NewMemoryContext()
	ctx.SetDefaultComment("Testing Default Comment")
	pkg := ctx.MustPackage("foo")
	pkg.File("bar")
	ctx.MustExport()

	expected := "// Testing Default Comment"
	assert.Equal(t, expected, ctx.Last()[:len(expected)])
}
