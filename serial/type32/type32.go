package type32

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/serial/rye"
)

type ErrTypeNotFound struct {
	reflect.Type
}

// Error fulfills error.
func (err ErrTypeNotFound) Error() string {
	return fmt.Sprintf("Type %s was not found", err.Type)
}

func checkLn(b []byte) []byte {
	// TODO: move this to slice
	ln := len(b)
	if cap(b) < ln+4 {
		cp := make([]byte, ln, ln*2+4)
		copy(cp, b)
		b = cp
	}
	return b
}

func sliceToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return rye.Deserialize.Uint32(b)
}
