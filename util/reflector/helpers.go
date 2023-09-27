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
