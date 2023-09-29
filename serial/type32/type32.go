package type32

import (
	"fmt"
	"reflect"
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
