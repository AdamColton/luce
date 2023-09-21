package reflector

import "reflect"

// Type creates a reflect.Type from the generic type without allocating memory.
// This is a wrapper around reflect.TypeOf((*T)(nil)).Elem().
func Type[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
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
