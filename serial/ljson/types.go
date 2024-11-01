package ljson

import (
	"reflect"
	"unsafe"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled.
type TypesContext struct {
	marshalers map[uintptr]unsafeMarshal
}

// NewTypesContext creates a TypesContext
func NewTypesContext() *TypesContext {
	return &TypesContext{
		marshalers: make(map[uintptr]unsafeMarshal),
	}
}

func (tctx *TypesContext) lazyGetter(t reflect.Type, self *unsafeMarshal) {
	tab := getITab(t)
	m, found := tctx.marshalers[tab]
	if found {
		*self = m
		return
	}

	*self = func(ptr unsafe.Pointer, ctx *MarshalContext) WriteNode {
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
