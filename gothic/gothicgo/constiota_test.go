package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestConstIotaBlock(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("foo")
	c := file.MustConstIotaBlock(IntType, "apple", "bannana", "cantaloupe", "date", "elderberry")
	c.Comment = "List of fruit"

	c = file.MustConstIotaBlock(ByteType, "read", "write", "execute")
	c.Iota = "1 << iota"

	ctx.MustExport()

	str := ctx.Last.String()
	assert.Contains(t, str, "// List of fruit\nconst (")
	assert.Contains(t, str, "apple int = iota")
	assert.Contains(t, str, "bannana\n")
	assert.Contains(t, str, "elderberry\n)")
	assert.Contains(t, str, "read byte = 1 << iota")
}
