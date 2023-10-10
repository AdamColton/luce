package reflector

import "reflect"

// Type creates a reflect.Type from the generic type without allocating memory.
// This is a wrapper around return reflect.TypeOf([0]T{}).Elem().
func Type[T any]() reflect.Type {
	return reflect.TypeOf([0]T{}).Elem()
}

// ToType returns reflect.Type unless it is already an instance reflect.Type.
func ToType(i any) reflect.Type {
	if t, ok := i.(reflect.Type); ok {
		return t
	}
	return reflect.TypeOf(i)
}

// ToValue returns reflect.Value of i unless it is already an instance of
// reflect.Value.
func ToValue(i any) reflect.Value {
	if v, ok := i.(reflect.Value); ok {
		return v
	}
	return reflect.ValueOf(i)
}

// ReturnsErrCheck checks the return values from a function call to see if the
// last value is an error.
func ReturnsErrCheck(returnVals []reflect.Value) error {
	if l := len(returnVals); l > 0 {
		err, ok := returnVals[l-1].Interface().(error)
		if ok {
			return err
		}
	}
	return nil
}

// CanNil reports wether k is a nilable kind.
func CanNil(k reflect.Kind) bool {
	return k == reflect.Chan ||
		k == reflect.Func ||
		k == reflect.Interface ||
		k == reflect.Map ||
		k == reflect.Pointer ||
		k == reflect.Slice
}

// IsNil reports whether its argument t is nil. Unlike the underlying t.IsNil,
// it will not panic.
func IsNil(t reflect.Value) bool {
	if CanNil(t.Kind()) {
		return t.IsNil()
	}
	return false
}

// Make a reflect.Value of type t.
func Make(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		return reflect.New(t)
	}
	return reflect.New(t).Elem()
}
