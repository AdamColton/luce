package reflector

import "reflect"

// Type returns the reflect.Type of T without allocating memory. It wraps
// reflect.TypeOf([0]T{}).Elem(), so unlike reflect.TypeOf(v) it also works for
// interface types such as error.
func Type[T any]() reflect.Type {
	// TODO: use reflect.TypeFor
	return reflect.TypeOf([0]T{}).Elem()
}
