package gothicgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExternalFunc(t *testing.T) {
	fn := MustPackageRef("foo").NewFuncRef("Foo", IntType.Unnamed())
	assert.Equal(t, "foo.Foo(x)", fn.Call(DefaultPrefixer, "x"))
}
