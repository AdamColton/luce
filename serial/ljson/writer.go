package ljson

import (
	"bytes"
	"io"
	"reflect"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/reflector"
)

// WriteContext is passed into a WriteNode. EscapeHtml makes strings escape <, >
// and &, and the SumWriter is where the json is written.
type WriteContext struct {
	EscapeHtml bool
	*luceio.SumWriter
	Nl, Tab string
	indent  string
}

// WriteNode writes a node of the json document. A WriteNode reports an error by
// setting the Err of its SumWriter.
type WriteNode func(ctx *WriteContext)

// String invokes the WriteNode and returns the data written as a string.
func (wn WriteNode) String() string {
	buf := bytes.NewBuffer(nil)
	wn.WriteTo(buf)
	return buf.String()
}

// WriteTo fulfills io.WriterTo and writes the WriteNode to the Writer.
func (wn WriteNode) WriteTo(w io.Writer) (int64, error) {
	wctx := &WriteContext{
		SumWriter: luceio.NewSumWriter(w),
	}
	wn(wctx)
	return wctx.Rets()
}

// Stringify marshals the value given and returns a json string. It returns the
// error from Marshal.
func Stringify[T, Ctx any](v T, ctx *MarshalContext[Ctx]) (string, error) {
	wn, err := Marshal(v, ctx)
	if err != nil {
		return "", err
	}
	return wn.String(), nil
}

// Export returns the names and types of the fields of the struct T that ctx
// writes: omitted fields are left out, generated fields are added and
// conditional fields are included only when their condition holds for ctx.
func Export[T, Ctx any](ctx *MarshalContext[Ctx]) (map[string]reflect.Type, error) {
	t := reflector.Type[T]()
	return ExportType(t, ctx)
}

// ExportType is Export for a reflect.Type, which can be a struct or a pointer to
// a struct. It returns an error for any other type.
func ExportType[Ctx any](t reflect.Type, ctx *MarshalContext[Ctx]) (map[string]reflect.Type, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, lerr.Str("expected struct or ptr to struct")
	}

	var vm valMarshaler[Ctx]
	ctx.TypesContext.get(t, &vm)
	if dg, ok := vm.(deferGet[Ctx]); ok {
		dg.get(ctx.TypesContext)
		vm = *dg.self
	}
	sm, ok := vm.(structMarshaler[Ctx])
	if !ok {
		sm = ctx.TypesContext.buildStructMarshal(t)
	}
	return sm.export(ctx), nil
}

type floodExport[Ctx any] struct {
	ctx *MarshalContext[Ctx]
}

func (fe floodExport[Ctx]) floodProc(t reflect.Type, add func(reflect.Type)) {
	k := t.Kind()
	if k == reflect.Array || k == reflect.Slice || k == reflect.Ptr {
		add(t.Elem())
	} else if k == reflect.Map {
		add(t.Elem())
		add(t.Key())
	} else if k == reflect.Struct {
		got := lerr.Must(ExportType(t, fe.ctx))
		for _, gt := range got {
			add(gt)
		}
	}
}

// ExportAll returns the exported fields of T and of every struct that can be
// reached from it through fields, pointers, slices, arrays and the keys and
// values of maps.
func ExportAll[T, Ctx any](ctx *MarshalContext[Ctx]) (map[reflect.Type]map[string]reflect.Type, error) {
	s := lset.New(reflector.Type[T]())
	fe := floodExport[Ctx]{ctx}
	s.Flood(fe.floodProc)

	structs, _ := filter.IsKind(reflect.Struct).SliceInPlace(s.Slice(nil))
	out := make(map[reflect.Type]map[string]reflect.Type, len(structs))
	for _, t := range structs {
		out[t] = lerr.Must(ExportType(t, ctx))
	}
	return out, nil
}
