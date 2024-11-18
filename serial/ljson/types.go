package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled.
type TypesContext struct {
	marshalers map[reflect.Type]valMarshaler
}

// NewTypesContext creates a TypesContext
func NewTypesContext() *TypesContext {
	return &TypesContext{
		marshalers: make(map[reflect.Type]valMarshaler),
	}
}

func (tctx *TypesContext) get(t reflect.Type, self *valMarshaler) {
	m, found := tctx.marshalers[t]
	if found {
		*self = m
		return
	}

	*self = func(v reflect.Value, ctx *MarshalContext) WriteNode {
		tctx := ctx.TypesContext

		m, found := tctx.marshalers[t]
		if !found {
			m = tctx.buildValueMarshaler(t)
			tctx.marshalers[t] = m
		}
		*self = m

		return (*self)(v, ctx)
	}
}

func (tctx *TypesContext) buildValueMarshaler(t reflect.Type) (m valMarshaler) {
	switch t.Kind() {
	case reflect.String:
		m = valMarshal(MarshalString)
	default:
		panic(lerr.Str("could not marshal " + t.String()))
	}
	return
}

func valMarshal[T any](m Marshaler[T]) valMarshaler {
	return func(v reflect.Value, ctx *MarshalContext) WriteNode {
		t := v.Interface().(T)
		return lerr.Must(m(t, ctx))
	}
}
