package type32

import (
	"github.com/adamcolton/luce/lerr"
)

// Sentinal Errors
const (
	ErrTooShort      = lerr.Str("TypeID32 too short")
	ErrNotRegistered = lerr.Str("No type registered")
	ErrSerNotT32     = lerr.Str("Serialize requires interface to be TypeIDer32")
	ErrNilZero       = lerr.Str("TypeID32Deserializer.Register) cannot register nil interface")
)

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen becuase
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
}
