package filter

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/util/reflector"
)

// TODO
type Type struct {
	Filter[reflect.Type]
}

// TODO
func NumInEq(n int) Type {
	return NumIn(EQ(n))
}

// TODO
func NumOutEq(n int) Type {
	return NumOut(EQ(n))
}

// TODO
func InType(n int, t reflect.Type) Type {
	return IsType(t).In(n)
}

// TODO
func OutType(n int, t reflect.Type) Type {
	return IsType(t).Out(n)
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

// Out checks the filter type's agument number i against the given filter.
func (t Type) Out(i int) Type {
	return Type{func(t2 reflect.Type) bool {
		nOut := t2.NumOut()
		idx := i
		if idx < 0 {
			idx += nOut
		}
		return t2 != nil && t2.Kind() == reflect.Func && idx >= 0 && idx < nOut && t.Filter(t2.Out(idx))
	}}
}

// TODO
func (t Type) Method() Filter[*reflector.Method] {
	return func(m *reflector.Method) bool {
		return t.Filter(m.Func.Type())
	}
}

// TODO
func (t Type) And(t2 Type) Type {
	return Type{t.Filter.And(t2.Filter)}
}

// TODO
func (t Type) Or(t2 Type) Type {
	return Type{t.Filter.Or(t2.Filter)}
}

// TODO
func (t Type) Not() Type {
	return Type{t.Filter.Not()}
}

// TypeChecker checks a value's type against a filter. It returns the underlying
// type. It returns an error if the type fails the underlying filter.
type TypeChecker func(i any) (reflect.Type, error)

// TODO
func (t Type) Check(errFn func(reflect.Type) error) TypeChecker {
	return func(i any) (reflect.Type, error) {
		it := reflector.ToType(i)
		if t.Filter(it) {
			return it, nil
		}
		return it, errFn(it)
	}
}

// Panic if i fails the underlying filter. Return the type of i if it succeeds.
func (tc TypeChecker) Panic(i any) reflect.Type {
	t, err := tc(i)
	if err != nil {
		panic(err)
	}
	return t
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

// NumOut checks the filter type's NumOut value against the given filter.
func NumOut(f Filter[int]) Type {
	return Type{func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumOut())
	}}
}

// TODO
func MethodName(f func(string) bool) Filter[*reflector.Method] {
	return func(m *reflector.Method) bool {
		return f(m.Name)
	}
}

func TypeErr(format string) func(t reflect.Type) error {
	return func(t reflect.Type) error {
		return fmt.Errorf(format, t)
	}
}
