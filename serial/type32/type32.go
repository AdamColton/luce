package type32

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/serial/rye"
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

func sliceToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return rye.Deserialize.Uint32(b)
}
