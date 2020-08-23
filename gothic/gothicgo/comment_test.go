package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestComment(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("foo")

	file.NewComment("This is a test")
	ctx.SetCommentWidth(10)
	file.NewComment("The sun was shining on the sea")

	assert.Equal(t, ErrCommentWidth, ctx.SetCommentWidth(0))

	ctx.MustExport()

	str := ctx.Last.String()
	assert.Contains(t, str, "// This is a test")
	assert.Contains(t, str, "// The sun")
	assert.Contains(t, str, "// was")
	assert.Contains(t, str, "// shining")
	assert.Contains(t, str, "// on the")
	assert.Contains(t, str, "// sea")
}
