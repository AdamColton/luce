package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"

	"github.com/testify/assert"
)

func TestExternalType(t *testing.T) {
	ref := MustPackageRef("foo")
	bar := ref.NewTypeRef("Bar", nil)

	i := NewImports(nil)
	bar.RegisterImports(i)
	str := bufpool.MustWriterToString(i)
	assert.Contains(t, str, "foo")
}
