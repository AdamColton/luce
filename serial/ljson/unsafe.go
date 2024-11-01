package ljson

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/lerr"
)

type unsafeMarshal func(ptr unsafe.Pointer, ctx *MarshalContext) WriteNode

func makeUnsafe[T any](m Marshaler[T]) unsafeMarshal {
	return func(ptr unsafe.Pointer, ctx *MarshalContext) WriteNode {
		t := *(*T)(ptr)
		return lerr.Must(m(t, ctx))
	}
}

func (ctx *TypesContext) buildUnsafeMarshaler(t reflect.Type) (m unsafeMarshal) {
	switch t.Kind() {
	case reflect.String:
		m = makeUnsafe(MarshalString)
	default:
		panic(lerr.Str("could not marshal " + t.String()))
	}

	return
}

type iface struct {
	tab  uintptr // actually *runtime.itab
	data unsafe.Pointer
}

func toIface[T any](ptr *T) iface {
	return *(*iface)(unsafe.Pointer(ptr))
}

func getITab(t reflect.Type) uintptr {
	a := reflect.New(t).Elem().Interface()
	return uintptr(toIface(&a).tab)
}
