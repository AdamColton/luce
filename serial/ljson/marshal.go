package ljson

import (
	"unsafe"

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
	var um unsafeMarshal[Ctx]
	ctx.lazyGetter(reflector.Type[T](), &um)
	m = func(v T, ctx *MarshalContext[Ctx]) (wn WriteNode, err error) {
		defer lerr.Recover(func(e error) { err = e })
		return um(unsafe.Pointer(&v), ctx), nil
	}
	return
}

// Marshaler is a function for creating a WriteNode for a value.
type Marshaler[T, Ctx any] func(v T, ctx *MarshalContext[Ctx]) (WriteNode, error)

// MarshalContext holds the context for marshaling values into WriteNodes,
// including the underlying TypesContext. The Context field holds an arbitrary
// data type that will be available during the marshaling phase.
type MarshalContext[Ctx any] struct {
	Context      Ctx
	TypesContext *TypesContext[Ctx]
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
	tab := getITab(reflector.Type[T]())
	ctx.marshalers[tab] = makeUnsafe(m)
}
