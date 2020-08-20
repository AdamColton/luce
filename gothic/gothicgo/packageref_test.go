package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestPackageRef(t *testing.T) {
	r := MustPackageRef("foo/bar")

	assert.Equal(t, "bar", r.Name())
	assert.Equal(t, "foo/bar", r.ImportPath())
	assert.Equal(t, `"foo/bar"`, r.ImportSpec())
}

func TestBadPackageRef(t *testing.T) {
	pkg, err := NewPackageRef("bad package ref")
	assert.Nil(t, pkg)
	assert.Error(t, err)
}
