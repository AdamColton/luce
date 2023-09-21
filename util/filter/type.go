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

// Elem checks the filter type's Elem against the underlying filter.
func (t Type) Elem() Type {
	return Type{func(t2 reflect.Type) (out bool) {
		return CanElem(t2) && t.Filter(t2.Elem())
	}}
}

// In checks the filter type's agument number i against the given filter.
func (t Type) In(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		nIn := t2.NumIn()
		idx := i
		if idx < 0 {
			idx += nIn
		}
		return t2 != nil && t2.Kind() == reflect.Func && idx >= 0 && idx < nIn && t.Filter(t2.In(idx))

	}}
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

// CanElem returns true if it is safe to call Elem on t.
func CanElem(t reflect.Type) bool {
	if t == nil {
		return false
	}
	k := t.Kind()
	return k == reflect.Array ||
		k == reflect.Chan ||
		k == reflect.Map ||
		k == reflect.Pointer ||
		k == reflect.Slice
}

// NumIn checks the filter type's NumIn value against the given filter.
func NumIn(f Filter[int]) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumIn())
	}}
}
