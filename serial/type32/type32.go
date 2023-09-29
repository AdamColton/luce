package type32

import (
	"fmt"
	"reflect"
)

// ErrTypeNotFound is returned when a type map does not contain the given
// uint32.
type ErrTypeNotFound struct {
	reflect.Type
}

// Error fulfills error.
func (err ErrTypeNotFound) Error() string {
	return fmt.Sprintf("type32: type %s was not found", err.Type)
}

func checkLn(b []byte) []byte {
	// confirm that b has additional capacity of at least 4
	// TODO: move this to slice
	ln := len(b)
	if cap(b) < ln+4 {
		cp := make([]byte, ln, ln*2+4)
		copy(cp, b)
		b = cp
	}
	return b
}
