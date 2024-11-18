package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled.
type TypesContext[Ctx any] struct {
	marshalers map[reflect.Type]valMarshaler[Ctx]
}

// NewTypesContext creates a TypesContext
func NewTypesContext[Ctx any]() *TypesContext[Ctx] {
	return &TypesContext[Ctx]{
		marshalers: make(map[reflect.Type]valMarshaler[Ctx]),
	}
}

func (tctx *TypesContext[Ctx]) get(t reflect.Type, self *valMarshaler[Ctx]) {
	m, found := tctx.marshalers[t]
	if found {
		*self = m
		return
	}

	*self = func(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
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

func (tctx *TypesContext[Ctx]) buildValueMarshaler(t reflect.Type) (m valMarshaler[Ctx]) {
	switch t.Kind() {
	case reflect.String:
		m = valMarshal(MarshalString[Ctx])
	default:
		panic(lerr.Str("could not marshal " + t.String()))
	}
	return
}

func valMarshal[T, Ctx any](m Marshaler[T, Ctx]) valMarshaler[Ctx] {
	return func(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
		t := v.Interface().(T)
		return lerr.Must(m(t, ctx))
	}
}
