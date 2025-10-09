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

// WriteContext is passed into a WriteNode
type WriteContext struct {
	EscapeHtml bool
	*luceio.SumWriter
	Nl, Tab string
	indent  string
}

// WriteNode writes a node of the json document
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

// Stringify marshals the value given and returns a json string.
func Stringify[T, Ctx any](v T, ctx *MarshalContext[Ctx]) (string, error) {
	wn, err := Marshal(v, ctx)
	if err != nil {
		return "", err
	}
	return wn.String(), nil
}

func Export[T, Ctx any](ctx *MarshalContext[Ctx]) (map[string]reflect.Type, error) {
	t := reflector.Type[T]()
	return ExportType(t, ctx)
}

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

func ExportAll[T, Ctx any](ctx *MarshalContext[Ctx]) (map[reflect.Type]map[string]reflect.Type, error) {
	s := lset.New(reflector.Type[T]())
	fe := floodExport[Ctx]{ctx}
	s.Flood(fe.floodProc)

	structs, _ := filter.IsKind(reflect.Struct).SliceInPlace(s.Slice(nil))
	out := make(map[reflect.Type]map[string]reflect.Type, len(structs))
	for _, t := range structs {
		ex, err := ExportType(t, ctx)
		if err != nil {
			return nil, err
		}
		out[t] = ex
	}
	return out, nil
}
