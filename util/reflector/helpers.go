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
