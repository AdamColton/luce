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
