package reflector

import "reflect"

// Type returns the reflect.Type of T without allocating memory. It wraps
// reflect.TypeOf([0]T{}).Elem(), so unlike reflect.TypeOf(v) it also works for
// interface types such as error.
func Type[T any]() reflect.Type {
	// TODO: use reflect.TypeFor
	return reflect.TypeOf([0]T{}).Elem()
}

// ToType returns the reflect.Type of i. If i is already a reflect.Type, it is
// returned as it is.
func ToType(i any) reflect.Type {
	if t, ok := i.(reflect.Type); ok {
		return t
	}
	return reflect.TypeOf(i)
}

// ToValue returns the reflect.Value of i. If i is already a reflect.Value, it is
// returned as it is.
func ToValue(i any) reflect.Value {
	if v, ok := i.(reflect.Value); ok {
		return v
	}
	return reflect.ValueOf(i)
}

// ReturnsErrCheck checks the return values from a function call to see if the
// last value is an error. It returns that error, or nil if there are no values,
// the last one is not an error or the error is nil.
func ReturnsErrCheck(returnVals []reflect.Value) error {
	if l := len(returnVals); l > 0 {
		err, ok := returnVals[l-1].Interface().(error)
		if ok {
			return err
		}
	}
	return nil
}

// CanNil reports whether k is a nilable kind: Chan, Func, Interface, Map,
// Pointer or Slice.
func CanNil(k reflect.Kind) bool {
	return k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Pointer ||
		k == reflect.Slice
}

// CanElem returns true if it is safe to call reflect.Type.Elem on a Type of
// kind k.
func CanElem(k reflect.Kind) bool {
	return k == reflect.Array ||
		k == reflect.Chan ||
		k == reflect.Map ||
		k == reflect.Pointer ||
		k == reflect.Slice
}

// Elem returns t.Elem() and true. If t is nil or its kind has no element type,
// it returns nil and false instead of panicking.
func Elem(t reflect.Type) (out reflect.Type, ok bool) {
	if t == nil {
		return
	}
	ok = CanElem(t.Kind())
	if ok {
		out = t.Elem()
	}
	return
}

// IsNil reports whether its argument t is nil. Unlike the underlying t.IsNil,
// it will not panic: it returns false for a kind that CanNil reports can't be
// nil.
func IsNil(t reflect.Value) bool {
	if CanNil(t.Kind()) {
		return t.IsNil()
	}
	return false
}

// Make creates a new zero reflect.Value of type t. If t is a pointer type, it
// allocates the value it points to and returns the pointer. Otherwise the value
// returned is addressable, so it can be set.
func Make(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		return reflect.New(t)
	}
	return reflect.New(t).Elem()
}

// Set attempts to set the 'to' value on the target and returns a bool to
// indicate success or failure. Will not panic. If the types differ and 'to' is
// an interface, the value inside the interface is set.
func Set(target, to reflect.Value) (out bool) {
	defer func() {
		recover()
	}()
	if target.Type() != to.Type() {
		if to.Kind() == reflect.Interface {
			to = to.Elem()
		}
	}
	target.Set(to)
	out = true
	return
}

// EnsurePointer returns v if it is a pointer. If it is not, it returns a pointer
// to a copy of v, so changes made through the pointer do not affect the
// original.
func EnsurePointer(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer {
		v2 := reflect.New(v.Type())
		v2.Elem().Set(v)
		v = v2
	}
	return v
}
