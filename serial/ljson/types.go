package ljson

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/reflector"
)

// TypesContext is used to create marshalers. Once a marshaler for a type is
// created within a TypesContext it is reused everytime that type needs to be
// marshaled.
type TypesContext[Ctx any] struct {
	marshalers      map[reflect.Type]valMarshaler[Ctx]
	fieldMarshalers map[FieldKey]valFieldMarshaler[Ctx]
	circularGuard   *lset.Set[reflect.Type]
	fieldGenerators map[reflect.Type][]valFieldMarshaler[Ctx]
}

// NewTypesContext creates a TypesContext
func NewTypesContext[Ctx any]() *TypesContext[Ctx] {
	return &TypesContext[Ctx]{
		marshalers:      make(map[reflect.Type]valMarshaler[Ctx]),
		fieldMarshalers: make(map[FieldKey]valFieldMarshaler[Ctx]),
		circularGuard:   lset.New[reflect.Type](),
		fieldGenerators: make(map[reflect.Type][]valFieldMarshaler[Ctx]),
	}
}

type deferGet[Ctx any] struct {
	self *valMarshaler[Ctx]
	t    reflect.Type
}

func (dg deferGet[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	tctx := ctx.TypesContext
	if tctx.circularGuard.Contains(dg.t) {
		panic(lerr.Str("circular type reference"))
	}
	tctx.circularGuard.Add(dg.t)

	m, found := tctx.marshalers[dg.t]
	if !found {
		m = tctx.buildValueMarshaler(dg.t)
		tctx.marshalers[dg.t] = m
	}
	*(dg.self) = m
	tctx.circularGuard.Remove(dg.t)

	return (*dg.self).marshalVal(v, ctx)
}

func (tctx *TypesContext[Ctx]) get(t reflect.Type, self *valMarshaler[Ctx]) {
	m, found := tctx.marshalers[t]
	if found {
		*self = m
		return
	}
	*self = deferGet[Ctx]{
		self: self,
		t:    t,
	}
}

func (tctx *TypesContext[Ctx]) buildValueMarshaler(t reflect.Type) (m valMarshaler[Ctx]) {
	switch t.Kind() {
	case reflect.Pointer:
		m = newMarshalPointer(t, tctx)
	case reflect.String:
		m = valMarshal(MarshalString[Ctx])
	case reflect.Slice:
		m = newMarshalSlice(t, tctx)
	case reflect.Struct:
		m = tctx.buildStructMarshal(t)
	case reflect.Map:
		m = newMarshalMap(t, tctx)
	case reflect.Interface:
		m = marshalInterface[Ctx]{}
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
	return m
}

type marshalPointer[Ctx any] struct {
	em valMarshaler[Ctx]
}

func newMarshalPointer[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) (out marshalPointer[Ctx]) {
	ctx.get(t.Elem(), &(out.em))
	return
}

func (mp marshalPointer[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	return ctx.guardMarshal(v.Elem(), mp.em)
}

type sliceWriter []WriteNode

func (s sliceWriter) writer(ctx *WriteContext) {
	ctx.WriteRune('[')
	if len(s) > 0 {
		s[0](ctx)
		for _, wn := range s[1:] {
			ctx.WriteRune(',')
			wn(ctx)
		}
	}
	ctx.WriteRune(']')
}

type marshalSlice[Ctx any] struct {
	em valMarshaler[Ctx]
}

func newMarshalSlice[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) (out marshalSlice[Ctx]) {
	et := t.Elem()
	ctx.get(et, &(out.em))
	return
}

func (ms marshalSlice[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	ln := v.Len()
	out := make(sliceWriter, ln)
	for i := 0; i < ln; i++ {
		out[i] = ctx.guardMarshal(v.Index(i), ms.em)
	}
	return out.writer
}

type marshalMap[Ctx any] struct {
	km, vm valMarshaler[Ctx]
}

func newMarshalMap[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) (out marshalMap[Ctx]) {
	kt := t.Key()
	ctx.get(kt, &(out.km))
	vt := t.Elem()
	ctx.get(vt, &(out.vm))
	return
}

func (mm marshalMap[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	out := make(StructWriter, 0, v.Len())
	mi := v.MapRange()
	for mi.Next() {
		kwn := ctx.guardMarshal(mi.Key(), mm.km)
		vwn := ctx.guardMarshal(mi.Value(), mm.vm)
		out = append(out, FieldWriter{
			Key:   kwn,
			Value: vwn,
		})
	}

	if ctx.Sort {
		out.sort()
	}
	return out.WriteNode
}

type marshalConvert[From, To, Ctx any] struct {
	vmTo valMarshaler[Ctx]
	fn   func(from From, ctx *MarshalContext[Ctx]) To
}

func (mc marshalConvert[From, To, Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	from := v.Interface().(From)
	to := mc.fn(from, ctx)
	return mc.vmTo.marshalVal(reflect.ValueOf(to), ctx)
}

// Convert a type before it is marshaled.
func Convert[From, To, Ctx any](fn func(from From, ctx *MarshalContext[Ctx]) To, ctx *TypesContext[Ctx]) {
	c := marshalConvert[From, To, Ctx]{
		fn: fn,
	}
	ctx.get(reflector.Type[To](), &(c.vmTo))

	ctx.marshalers[reflector.Type[From]()] = c
}

type marshalInterface[Ctx any] struct{}

func (marshalInterface[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	t := v.Type()
	if t.Kind() == reflect.Interface {
		v = v.Elem()
		t = v.Type()
	}
	var em valMarshaler[Ctx]
	ctx.TypesContext.get(t, &em)
	return ctx.guardMarshal(v, em)
}
