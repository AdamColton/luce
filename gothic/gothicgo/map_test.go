package gothicgo

import (
	"testing"
)

func TestMapRegisterImports(t *testing.T) {
	i := NewImports(nil)
	mp := MapOf(IntType, StringType)

	mp.RegisterImports(i)

	// todo: after external type defs are done use types that will cause
	// registration.
}
