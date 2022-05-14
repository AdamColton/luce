package filter

import "reflect"

// IsKind filter checks Kind.
func IsKind(kind reflect.Kind) Filter[reflect.Type] {
	return func(t reflect.Type) bool {
		return t != nil && t.Kind() == kind
	}
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

// Elem checks the filter type's Elem against the underlying filter.
func Elem(f Filter[reflect.Type]) Filter[reflect.Type] {
	return func(t reflect.Type) (out bool) {
		return CanElem(t) && f(t.Elem())
	}
}

// IsType checks the filters Type against i. IsType will call reflect.TypeOf on
// i if it is not reflect.Type.
func IsType(i any) Filter[reflect.Type] {
	t1 := ToType(i)
	return func(t2 reflect.Type) bool {
		return t1 == t2
	}
}

// IsNilRef expects nilRefPtr to be a nil pointer. It checks the filter type
// against the type pointed to. This is useful for cases like interfaces where
// passing in a nil to IsType will fail.
func IsNilRef(nilRefPtr any) Filter[reflect.Type] {
	t1 := reflect.TypeOf(nilRefPtr).Elem()
	return func(t2 reflect.Type) bool {
		return t1 == t2
	}
}

// NumIn checks the filter type's NumIn value against the given filter.
func NumIn(f Filter[int]) Filter[reflect.Type] {
	return func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && f(t.NumIn())
	}
}

// In checks the filter type's agument number i against the given filter.
func In(i int, f Filter[reflect.Type]) Filter[reflect.Type] {
	return func(t reflect.Type) (out bool) {
		return t != nil && t.Kind() == reflect.Func && i < t.NumIn() && f(t.In(i))
	}
}

// ToType returns TypeOf unless it is already a Type.
func ToType(i any) reflect.Type {
	if t, ok := i.(reflect.Type); ok {
		return t
	}
	return reflect.TypeOf(i)
}

// TypeChecker checks a value's type against a filter. It returns the underlying
// type. It returns an error if the type fails the underlying filter.
type TypeChecker func(i any) (reflect.Type, error)

// TypeCheck produces a TypeChecker. Given a value, it will get the type and if
// that type fails the provided filter, it will return err. It also returns the
// type to provide a method to get and check the type of a value simultaneously.
func TypeCheck(f Filter[reflect.Type], err error) TypeChecker {
	return func(i any) (reflect.Type, error) {
		t := ToType(i)
		if f(t) {
			return t, nil
		}
		return t, err
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
