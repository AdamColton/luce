package reflector

import (
	"reflect"
	"unsafe"
)

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

// CanElem returns true if it is safe to call Elem on k.
func CanElem(k reflect.Kind) bool {
	return k == reflect.Array ||
		k == reflect.Chan ||
		k == reflect.Map ||
		k == reflect.Pointer ||
		k == reflect.Slice
}

// Elem calls t.Elem if it is safe to do so.
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

// Set attempts to set the 'to' value on the target and returns a bool to
// indicate success or failure. Will not panic.
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

// EnsurePointer check the Kind of v and if it is not a pointer, returns a
// pointer to v.
func EnsurePointer(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Pointer {
		v2 := reflect.New(v.Type())
		v2.Elem().Set(v)
		v = v2
	}
	return v
}

// UnsafeByteSlice returns a byte slice holding the memory of an arbitrary type.
func UnsafeByteSlice[T any](t T) []byte {
	ln := unsafe.Sizeof(t)
	p := (*byte)(unsafe.Pointer(&t))
	return unsafe.Slice(p, ln)
}
