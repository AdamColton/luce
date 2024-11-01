package ljson

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/lerr"
)

type unsafeMarshal[Ctx any] func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode

func makeUnsafe[T, Ctx any](m Marshaler[T, Ctx]) unsafeMarshal[Ctx] {
	return func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		t := *(*T)(ptr)
		return lerr.Must(m(t, ctx))
	}
}

func (ctx *TypesContext[Ctx]) buildUnsafeMarshaler(t reflect.Type) (m unsafeMarshal[Ctx]) {
	switch t.Kind() {
	case reflect.String:
		m = makeUnsafe(MarshalString[Ctx])
	case reflect.Struct:
		m = ctx.buildStructMarsal(t).unsafeMarshal
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
