package ltype

import (
	"reflect"

	"github.com/adamcolton/luce/util/reflector"
)

// All base types are defined in ltype
var (
	Err       = reflector.Type[error]()
	String    = reflector.Type[string]()
	Byte      = reflector.Type[byte]()
	Bool      = reflector.Type[bool]()
	ByteSlice = reflect.SliceOf(Byte)

	Int   = reflector.Type[int]()
	Int8  = reflector.Type[int8]()
	Int16 = reflector.Type[int16]()
	Int32 = reflector.Type[int32]()
	Int64 = reflector.Type[int64]()

	Uint   = reflector.Type[uint]()
	Uint8  = reflector.Type[uint8]()
	Uint16 = reflector.Type[uint16]()
	Uint32 = reflector.Type[uint32]()
	Uint64 = reflector.Type[uint64]()

	Float32 = reflector.Type[float32]()
	Float64 = reflector.Type[float64]()
)

// CheckStructPtr returns nil if t is not a pointer to a struct. If it is, it
// returns the struct.
func CheckStructPtr(t reflect.Type) reflect.Type {
	if t == nil || t.Kind() != reflect.Ptr {
		return nil
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return nil
	}
	return t
}
