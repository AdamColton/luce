package gothicgo_test

import (
	"testing"

	"github.com/adamcolton/luce/gothic/gothicgo"
	"github.com/testify/assert"
)

func TestDefaultPrefixer(t *testing.T) {
	bar := gothicgo.MustExternalPackageRef("foo/bar")
	prefix := gothicgo.DefaultPrefixer.Prefix(bar)
	assert.Equal(t, "bar.", prefix)

	prefix = gothicgo.DefaultPrefixer.Prefix(gothicgo.PkgBuiltin())
	assert.Equal(t, "", prefix)
}
