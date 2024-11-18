package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled.
type TypesContext[Ctx any] struct {
	marshalers      map[reflect.Type]valMarshaler[Ctx]
	fieldMarshalers map[FieldKey]valFieldMarshaler[Ctx]
}

// NewTypesContext creates a TypesContext
func NewTypesContext[Ctx any]() *TypesContext[Ctx] {
	return &TypesContext[Ctx]{
		marshalers:      make(map[reflect.Type]valMarshaler[Ctx]),
		fieldMarshalers: make(map[FieldKey]valFieldMarshaler[Ctx]),
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
	case reflect.Struct:
		m = tctx.buildStructMarshal(t).valMarshal
	case reflect.Int:
		m = valMarshal(MarshalInt[int, Ctx])
	case reflect.Int8:
		m = valMarshal(MarshalInt[int8, Ctx])
	case reflect.Int16:
		m = valMarshal(MarshalInt[int16, Ctx])
	case reflect.Int32:
		m = valMarshal(MarshalInt[int32, Ctx])
	case reflect.Int64:
		m = valMarshal(MarshalInt[int64, Ctx])
	case reflect.Uint:
		m = valMarshal(MarshalUint[uint, Ctx])
	case reflect.Uint8:
		m = valMarshal(MarshalUint[uint8, Ctx])
	case reflect.Uint16:
		m = valMarshal(MarshalUint[uint16, Ctx])
	case reflect.Uint32:
		m = valMarshal(MarshalUint[uint32, Ctx])
	case reflect.Uint64:
		m = valMarshal(MarshalUint[uint64, Ctx])
	case reflect.Bool:
		m = valMarshal(MarshalBool[Ctx])
	case reflect.Float64:
		m = valMarshal(MarshalFloat[float64, Ctx])
	case reflect.Float32:
		m = valMarshal(MarshalFloat[float32, Ctx])
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
