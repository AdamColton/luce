package type32

import (
	"github.com/adamcolton/luce/lerr"
)

// Sentinel errors
const (
	// ErrTooShort is returned by GetType when the data is shorter than the 4
	// bytes of a type ID.
	ErrTooShort = lerr.Str("TypeID32 too short")
	// ErrNotRegistered is returned by GetType when the type ID is not in the
	// TypeMap.
	ErrNotRegistered = lerr.Str("No type registered")
	// ErrSerNotT32 is the error for serializing a value that does not fulfill
	// TypeIDer32.
	ErrSerNotT32 = lerr.Str("Serialize requires interface to be TypeIDer32")
	// ErrNilZero is returned by RegisterType when it is given nil.
	ErrNilZero = lerr.Str("TypeID32Deserializer.Register) cannot register nil interface")
)

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen because
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
}
