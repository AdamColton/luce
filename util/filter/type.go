package filter

import (
	"reflect"

	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/util/reflector"
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

// Elem checks the filter type's Elem against the underlying filter.
func (t Type) Elem() Type {
	return Type{func(t2 reflect.Type) (out bool) {
		e, ok := reflector.Elem(t2)
		return ok && t.Filter(e)
	}}
}

// In checks the filter type's agument number i against the given filter. If i
// is less than 0, it will be indexed relative to the number of arguments, so -1
// will return the last argument.
func (t Type) In(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		idx, inRange := ints.Idx(i, t2.NumIn())
		return inRange && t2 != nil && t2.Kind() == reflect.Func && t.Filter(t2.In(idx))

	}}
}

// Out checks the filter type's agument number i against the given filter. If i
// is less than 0, it will be indexed relative to the number of returns, so -1
// will return the last return.
func (t Type) Out(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		idx, inRange := ints.Idx(i, t2.NumOut())
		return inRange && t2 != nil && t2.Kind() == reflect.Func && t.Filter(t2.Out(idx))
	}}
}

// TypeChecker checks a value's type against a filter. It returns the underlying
// type. It returns an error if the type fails the underlying filter.
type TypeChecker func(i any) (reflect.Type, error)

// Check creates a TypeChecker from a Type filter. It uses reflector.ToType,
// so that it can accept either a reflect.Type and use it directly or an
// interface which it will call reflect.ToType on.
func (t Type) Check(errFn func(reflect.Type) error) TypeChecker {
	return func(i any) (reflect.Type, error) {
		it := reflector.ToType(i)
		if t.Filter(it) {
			return it, nil
		}
		return it, errFn(it)
	}
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

// NumIn checks the filter type's NumIn value against the given filter.
func NumIn(f Filter[int]) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumIn())
	}}
}

// NumOut checks the filter type's NumOut value against the given filter.
func NumOut(f Filter[int]) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumOut())
	}}
}
