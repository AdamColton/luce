package ljson

import (
	"bytes"
	"reflect"

	"github.com/adamcolton/luce/ds/lset"
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
		return vm.marshalVal(reflect.ValueOf(v), ctx), nil
	}
	return
}

type valMarshaler[Ctx any] interface {
	marshalVal(reflect.Value, *MarshalContext[Ctx]) WriteNode
}

// Marshaler is a function for creating a WriteNode for a value.
type Marshaler[T, Ctx any] func(v T, ctx *MarshalContext[Ctx]) (WriteNode, error)

func (m Marshaler[T, Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	t := v.Interface().(T)
	return lerr.Must(m(t, ctx))
}

// MarshalContext holds the context for marshaling values into WriteNodes,
// including the underlying TypesContext. The Context field holds an arbitrary
// data type that will be available during the marshaling phase.
type MarshalContext[Ctx any] struct {
	Context       Ctx
	TypesContext  *TypesContext[Ctx]
	circularGuard *lset.Set[reflect.Value]

	// Setting Sort to true will sort the keys on structs and maps.
	// This is useful for testing because it produces consistent output
	// but may be skipped for effiency in production.
	Sort bool
}

// NewMarshalContext creates a MarshalContext using the TypesContext.
func (tctx *TypesContext[Ctx]) NewMarshalContext(ctx Ctx) *MarshalContext[Ctx] {
	return &MarshalContext[Ctx]{
		Context:       ctx,
		TypesContext:  tctx,
		circularGuard: lset.New[reflect.Value](),
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

func (ctx *MarshalContext[Ctx]) initGuard(v reflect.Value) {
	if ctx.circularGuard.Contains(v) {
		panic(lerr.Str("marshaled object contains circular reference"))
	}
	ctx.circularGuard.Add(v)
}

func (ctx *MarshalContext[Ctx]) guardMarshal(v reflect.Value, fn valMarshaler[Ctx]) WriteNode {
	ctx.initGuard(v)
	wn := fn.marshalVal(v, ctx)
	ctx.circularGuard.Remove(v)
	return wn
}

func (ctx *MarshalContext[Ctx]) guardFieldMarshaler(name string, v reflect.Value, fm valFieldMarshaler[Ctx]) (string, WriteNode) {
	ctx.initGuard(v)
	name, wn := fm(name, v, ctx)
	ctx.circularGuard.Remove(v)
	return name, wn
}

// Serialize fulfills serial.Serializer.
func (ctx *MarshalContext[Ctx]) Serialize(i any, buf []byte) (data []byte, err error) {
	defer lerr.Recover(func(e error) { err = e })

	v := reflect.ValueOf(i)
	var um valMarshaler[Ctx]
	ctx.TypesContext.get(v.Type(), &um)

	out := bytes.NewBuffer(buf)
	_, err = um.marshalVal(v, ctx).WriteTo(out)
	data = out.Bytes()
	return
}
