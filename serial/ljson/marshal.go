package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

// Marshal a value into a WriteNode
func Marshal[T any](v T, ctx *MarshalContext) (wn WriteNode, err error) {
	m, err := getMarshaler[T](ctx.TypesContext)
	if err != nil {
		return nil, err
	}
	wn, err = m(v, ctx)
	return
}

func getMarshaler[T any](ctx *TypesContext) (m Marshaler[T], err error) {
	defer lerr.Recover(func(e error) { err = e })
	var vm valMarshaler
	ctx.get(reflector.Type[T](), &vm)
	m = func(v T, ctx *MarshalContext) (wn WriteNode, err error) {
		defer lerr.Recover(func(e error) { err = e })
		return vm(reflect.ValueOf(v), ctx), nil
	}
	return
}

type valMarshaler func(reflect.Value, *MarshalContext) WriteNode

// Marshaler is a function for creating a WriteNode for a value.
type Marshaler[T any] func(v T, ctx *MarshalContext) (WriteNode, error)

// MarshalContext holds the context for marshaling values into WriteNodes,
// including the underlying TypesContext
type MarshalContext struct {
	TypesContext *TypesContext
}

// NewMarshalContext creates a MarshalContext using the TypesContext.
func (tctx *TypesContext) NewMarshalContext() *MarshalContext {
	return &MarshalContext{
		TypesContext: tctx,
	}
}

// NewMarshalContext creates both a NewTypesContext and a NewMarshalContext.
func NewMarshalContext() *MarshalContext {
	return NewTypesContext().NewMarshalContext()
}