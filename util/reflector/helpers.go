package reflector

import "reflect"

// Type creates a reflect.Type from the generic type without allocating memory.
// This is a wrapper around return reflect.TypeOf([0]T{}).Elem().
func Type[T any]() reflect.Type {
	return reflect.TypeOf([0]T{}).Elem()
}
