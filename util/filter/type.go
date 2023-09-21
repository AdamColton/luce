package filter

import (
	"reflect"
)

// Type is a wrapper around Filter[reflect.Type] to provide helper logic
// for type filtering.
type Type struct {
	Filter[reflect.Type]
}

// OnInterface applies the filter to the TypeOf i.
func (t Type) OnInterface(i any) bool {
	return t.Filter(reflect.TypeOf(i))
}

// IsKind creates a Type filter that returns true when given a type that
// matches the specified kind.
func IsKind(kind reflect.Kind) Type {
	return Type{func(t reflect.Type) bool {
		return t != nil && t.Kind() == kind
	}}
}

// IsType creates a filter using referenceType. Returns true if the filterType
// is the same as referenceType.
func IsType(referenceType reflect.Type) Type {
	return Type{func(filterType reflect.Type) bool {
		return referenceType == filterType
	}}
}
