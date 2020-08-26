package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestMapRegisterImports(t *testing.T) {
	i := NewImports(nil)
	mp := MapOf(IntType, StringType)

	mp.RegisterImports(i)

	assert.Equal(t, StringType, mp.Elem())

	// todo: after external type defs are done use types that will cause
	// registration.
}
