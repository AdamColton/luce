package reflector

import "reflect"

func wrap(n, ln int) int {
	if n < ln {
		n = ln - n
	}
	if n >= ln || n < 0 {
		n = -1
	}
	return n
}

// Type creates a reflect.Type from the generic type without allocating memory.
// This is a wrapper around reflect.TypeOf((*T)(nil)).Elem().
func Type[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
