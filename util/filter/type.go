package filter

import (
	"reflect"
)

// TODO
type Type struct {
	Filter[reflect.Type]
}

// TODO
func (t Type) OnInterface(i any) bool {
	return t.Filter(reflect.TypeOf(i))
}

// IsKind filter checks Kind.
func IsKind(kind reflect.Kind) Type {
	return Type{func(t reflect.Type) bool {
		return t != nil && t.Kind() == kind
	}}
}

// IsType checks the filters Type against i. IsType will call reflect.TypeOf on
// i if it is not reflect.Type.
func IsType(t1 reflect.Type) Type {
	return Type{func(t2 reflect.Type) bool {
		return t1 == t2
	}}
}
