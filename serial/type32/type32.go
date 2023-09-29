package type32

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/serial/rye"
)

// ErrTypeNotFound is returned when a prefixer has no uint32 for a type: a value
// that does not fulfill TypeIDer32, or a type that is not in the map.
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
