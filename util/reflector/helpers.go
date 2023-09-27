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
