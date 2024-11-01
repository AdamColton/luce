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
	case reflect.Int:
		m = makeUnsafe(MarshalInt[int, Ctx])
	case reflect.Int8:
		m = makeUnsafe(MarshalInt[int8, Ctx])
	case reflect.Int16:
		m = makeUnsafe(MarshalInt[int16, Ctx])
	case reflect.Int32:
		m = makeUnsafe(MarshalInt[int32, Ctx])
	case reflect.Int64:
		m = makeUnsafe(MarshalInt[int64, Ctx])
	case reflect.Uint:
		m = makeUnsafe(MarshalUint[uint, Ctx])
	case reflect.Uint8:
		m = makeUnsafe(MarshalUint[uint8, Ctx])
	case reflect.Uint16:
		m = makeUnsafe(MarshalUint[uint16, Ctx])
	case reflect.Uint32:
		m = makeUnsafe(MarshalUint[uint32, Ctx])
	case reflect.Uint64:
		m = makeUnsafe(MarshalUint[uint64, Ctx])
	case reflect.Bool:
		m = makeUnsafe(MarshalBool[Ctx])
	case reflect.Float64:
		m = makeUnsafe(MarshalFloat[float64, Ctx])
	case reflect.Float32:
		m = makeUnsafe(MarshalFloat[float32, Ctx])
	case reflect.Pointer:
		m = marshalPointer(t, ctx)
	case reflect.Slice:
		m = unsafeSliceMarshal(t, ctx)
	case reflect.Map:
		m = unsafeMapMarshal(t, ctx)
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

func marshalPointer[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) unsafeMarshal[Ctx] {
	var em unsafeMarshal[Ctx]
	ctx.lazyGetter(t.Elem(), &em)
	return func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		ptAt := *(*unsafe.Pointer)(ptr)
		return em(ptAt, ctx)
	}
}

type sliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}

type sliceWriter []WriteNode

func unsafeSliceMarshal[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) (m unsafeMarshal[Ctx]) {
	et := t.Elem()
	var em unsafeMarshal[Ctx]
	ctx.lazyGetter(et, &em)
	size := et.Size()
	return func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		s := *(*sliceHeader)(ptr)
		end := uintptr(s.Len)*size + s.Data
		out := make(sliceWriter, 0, s.Len)
		for ptr := s.Data; ptr < end; ptr += size {
			wn := ctx.guardMarshal(et, unsafe.Pointer(ptr), em)
			if wn != nil {
				out = append(out, wn)
			}
		}
		return out.writer
	}
}

func (s sliceWriter) writer(ctx *WriteContext) {
	ctx.WriteRune('[')
	for i, wn := range s {
		if i > 0 {
			ctx.WriteRune(',')
		}
		wn(ctx)
	}
	ctx.WriteRune(']')
}

func unsafeMapMarshal[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) (m unsafeMarshal[Ctx]) {
	kt := t.Key()
	var km unsafeMarshal[Ctx]
	ctx.lazyGetter(kt, &km)
	vt := t.Elem()

	// I don't entirely understand this.
	// I guess an interface always holds the direct object
	for vt.Kind() == reflect.Pointer {
		vt = vt.Elem()
	}
	var vm unsafeMarshal[Ctx]
	ctx.lazyGetter(vt, &vm)
	return func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		v := reflect.NewAt(t, ptr).Elem()
		out := make(StructWriter, 0, v.Len())
		mi := v.MapRange()
		for mi.Next() {
			k := mi.Key().Interface()
			ki := toIface(&k)
			kwn := ctx.guardMarshal(kt, ki.data, km)

			v := mi.Value().Interface()
			vi := toIface(&v)
			vwn := ctx.guardMarshal(vt, vi.data, vm)

			out = append(out, FieldWriter{
				Key:   kwn,
				Value: vwn,
			})
		}

		if ctx.Sort {
			out.sort()
		}
		return out.WriteNode
	}
}
