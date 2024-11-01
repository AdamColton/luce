package ljson

import (
	"reflect"
	"unsafe"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled. Ctx defines the type for the Context field during the marshaling
// phase.
type TypesContext[Ctx any] struct {
	marshalers map[uintptr]unsafeMarshal[Ctx]
}

// NewTypesContext creates a TypesContext
func NewTypesContext[Ctx any]() *TypesContext[Ctx] {
	return &TypesContext[Ctx]{
		marshalers: make(map[uintptr]unsafeMarshal[Ctx]),
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
		m, found := tctx.marshalers[tab]
		if !found {
			m = tctx.buildUnsafeMarshaler(t)
			tctx.marshalers[tab] = m
		}
		*self = m

		return (*self)(ptr, ctx)
	}
}
