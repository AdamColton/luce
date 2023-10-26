package ltype

import (
	"reflect"

	"github.com/adamcolton/luce/util/reflector"
)

var (
	Err       = reflector.Type[error]()
	String    = reflector.Type[string]()
	Byte      = reflector.Type[byte]()
	ByteSlice = reflect.SliceOf(Byte)

	Int   = reflector.Type[int]()
	Int8  = reflector.Type[int8]()
	Int16 = reflector.Type[int16]()
	Int32 = reflector.Type[int32]()
	Int64 = reflector.Type[int64]()

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
