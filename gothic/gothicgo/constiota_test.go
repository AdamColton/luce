package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestConstIotaBlock(t *testing.T) {
	ctx, file := newFile("foo")

	c := file.MustConstIotaBlock(IntType, "apple", "bannana", "cantaloupe", "date", "elderberry")
	c.Comment = "List of fruit"

	c = file.MustConstIotaBlock(ByteType, "read", "write", "execute")
	c.Iota = "1 << iota"

	ctx.MustExport()

	assert.Contains(t, ctx.Last(), "// List of fruit\nconst (")
	assert.Contains(t, ctx.Last(), "apple int = iota")
	assert.Contains(t, ctx.Last(), "bannana\n")
	assert.Contains(t, ctx.Last(), "elderberry\n)")
	assert.Contains(t, ctx.Last(), "read byte = 1 << iota")
}
