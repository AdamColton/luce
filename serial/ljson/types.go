package ljson

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled. Ctx defines the type for the Context field during the marshaling
// phase.
type TypesContext[Ctx any] struct {
	marshalers      map[uintptr]unsafeMarshal[Ctx]
	fieldMarshal    map[FieldKey]unsafeFieldMarshaler[Ctx]
	circularGuard   *lset.Set[reflect.Type]
	fieldGenerators map[reflect.Type][]unsafeFieldMarshaler[Ctx]
}

// NewTypesContext creates a TypesContext
func NewTypesContext[Ctx any]() *TypesContext[Ctx] {
	return &TypesContext[Ctx]{
		marshalers:      make(map[uintptr]unsafeMarshal[Ctx]),
		fieldMarshal:    make(map[FieldKey]unsafeFieldMarshaler[Ctx]),
		circularGuard:   lset.New[reflect.Type](),
		fieldGenerators: make(map[reflect.Type][]unsafeFieldMarshaler[Ctx]),
	}
}

func (tctx *TypesContext[Ctx]) lazyGetter(t reflect.Type, self *unsafeMarshal[Ctx]) {
	tab := getITab(t)
	m, found := tctx.marshalers[tab]
	if found {
		*self = m
		return
	}

	*self = func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		tctx := ctx.TypesContext
		if tctx.circularGuard.Contains(t) {
			panic(lerr.Str("circular type reference"))
		}
		tctx.circularGuard.Add(t)
		m, found := tctx.marshalers[tab]
		if !found {
			m = tctx.buildUnsafeMarshaler(t)
			tctx.marshalers[tab] = m
		}
		*self = m
		tctx.circularGuard.Remove(t)

		return (*self)(ptr, ctx)
	}
}

// Convert a type before it is marshaled.
func Convert[From, To, Ctx any](fn func(from From, ctx *MarshalContext[Ctx]) To, ctx *TypesContext[Ctx]) {
	var umTo unsafeMarshal[Ctx]
	ctx.lazyGetter(reflector.Type[To](), &umTo)
	umFrom := func(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
		from := *(*From)(ptr)
		to := fn(from, ctx)
		return umTo(unsafe.Pointer(&to), ctx)
	}
	tab := getITab(reflector.Type[From]())
	ctx.marshalers[tab] = umFrom
}
