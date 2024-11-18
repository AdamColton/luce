package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

// Marshal a value into a WriteNode
func Marshal[T, Ctx any](v T, ctx *MarshalContext[Ctx]) (wn WriteNode, err error) {
	m, err := getMarshaler[T](ctx.TypesContext)
	if err != nil {
		return nil, err
	}
	wn, err = m(v, ctx)
	return
}

func getMarshaler[T, Ctx any](ctx *TypesContext[Ctx]) (m Marshaler[T, Ctx], err error) {
	defer lerr.Recover(func(e error) { err = e })
	var vm valMarshaler[Ctx]
	ctx.get(reflector.Type[T](), &vm)
	m = func(v T, ctx *MarshalContext[Ctx]) (wn WriteNode, err error) {
		defer lerr.Recover(func(e error) { err = e })
		return vm(reflect.ValueOf(v), ctx), nil
	}
	return
}

type valMarshaler[Ctx any] func(reflect.Value, *MarshalContext[Ctx]) WriteNode

// Marshaler is a function for creating a WriteNode for a value.
type Marshaler[T, Ctx any] func(v T, ctx *MarshalContext[Ctx]) (WriteNode, error)

// MarshalContext holds the context for marshaling values into WriteNodes,
// including the underlying TypesContext. The Context field holds an arbitrary
// data type that will be available during the marshaling phase.
type MarshalContext[Ctx any] struct {
	Context      Ctx
	TypesContext *TypesContext[Ctx]

	// Setting Sort to true will sort the keys on structs and maps.
	// This is useful for testing because it produces consistent output
	// but may be skipped for effiency in production.
	Sort bool
}

// NewMarshalContext creates a MarshalContext using the TypesContext.
func (tctx *TypesContext[Ctx]) NewMarshalContext(ctx Ctx) *MarshalContext[Ctx] {
	return &MarshalContext[Ctx]{
		Context:      ctx,
		TypesContext: tctx,
	}
}

// NewMarshalContext creates both a NewTypesContext and a NewMarshalContext.
func NewMarshalContext[Ctx any](ctx Ctx) *MarshalContext[Ctx] {
	return NewTypesContext[Ctx]().NewMarshalContext(ctx)
}

// AddMarshaler to the TypesContext. This should be invoked before the
// TypesContext is used to marshal any values.
func AddMarshaler[T, Ctx any](m Marshaler[T, Ctx], ctx *TypesContext[Ctx]) {
	t := reflector.Type[T]()
	ctx.marshalers[t] = valMarshal(m)
}
