package ljson

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// FieldKey is used to identify a Field by name on a Type.
type FieldKey struct {
	reflect.Type
	Name string
}

type unsafeFieldMarshaler[Ctx any] func(name string, ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) (string, WriteNode)

func (tctx *TypesContext[Ctx]) lazyFieldMarshal(t reflect.Type, f reflect.StructField, self *unsafeFieldMarshaler[Ctx]) {
	key := FieldKey{
		Type: t,
		Name: f.Name,
	}
	fm, found := tctx.fieldMarshal[key]
	if found {
		*self = fm
		return
	}
	*self = func(name string, ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		fm, found := tctx.fieldMarshal[key]
		if !found {
			fm = defaultFieldMarshaler(f.Type, tctx)
			tctx.fieldMarshal[key] = fm
		}

		*self = fm
		return (*self)(name, ptr, ctx)
	}
}

func (tctx *TypesContext[Ctx]) buildStructMarsal(t reflect.Type) structMarshaler[Ctx] {
	return tctx.buildStructMarsalRecurse(t, 0, t.NumField(), 0)
}

func (tctx *TypesContext[Ctx]) buildStructMarsalRecurse(t reflect.Type, i, n, ln int) structMarshaler[Ctx] {
	if i == n {
		return make(structMarshaler[Ctx], ln)
	}
	f := t.Field(i)
	if !f.IsExported() {
		return tctx.buildStructMarsalRecurse(t, i+1, n, ln)
	}
	var fm unsafeFieldMarshaler[Ctx]
	tctx.lazyFieldMarshal(t, f, &fm)
	if fm == nil {
		return tctx.buildStructMarsalRecurse(t, i+1, n, ln)
	}
	sm := tctx.buildStructMarsalRecurse(t, i+1, n, ln+1)
	sm[ln] = fieldMarhsal[Ctx]{
		offset:               f.Offset,
		name:                 f.Name,
		unsafeFieldMarshaler: &fm,
		t:                    f.Type,
	}
	return sm
}

type fieldMarhsal[Ctx any] struct {
	offset uintptr
	name   string
	t      reflect.Type
	*unsafeFieldMarshaler[Ctx]
}

type structMarshaler[Ctx any] []fieldMarhsal[Ctx]

func (sm structMarshaler[Ctx]) unsafeMarshal(ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) WriteNode {
	start := uintptr(ptr)
	out := make(StructWriter, 0, len(sm))
	for _, fm := range sm {
		var fw FieldWriter

		fp := unsafe.Pointer(fm.offset + start)
		var name string
		name, fw.Value = (*fm.unsafeFieldMarshaler)(fm.name, fp, ctx)
		if name != "" {
			fw.Key = lerr.Must(MarshalString(name, ctx))
			out = append(out, fw)
		}
	}

	if ctx.Sort {
		out.sort()
	}
	return out.WriteNode
}

// FieldWriter writes the key and value of a single field in a json object.
type FieldWriter struct {
	Key, Value WriteNode
	sortKey    string
}

// StructWriter writes a json object.
type StructWriter []FieldWriter

func (s StructWriter) sort() {
	wctx := &WriteContext{}
	for i, fw := range s {
		buf, sw := luceio.BufferSumWriter()
		wctx.SumWriter = sw
		fw.Key(wctx)
		s[i].sortKey = buf.String()
	}
	slice.New(s).Sort(func(i, j FieldWriter) bool {
		return i.sortKey < j.sortKey
	})
}

// WriteNode method to actually fulfill WriteNode on StructWriter.
func (s StructWriter) WriteNode(ctx *WriteContext) {
	ctx.WriteRune('{')
	for i, fw := range s {
		if i > 0 {
			ctx.WriteRune(',')
		}
		fw.Key(ctx)
		ctx.WriteStrings(":")
		fw.Value(ctx)
	}
	ctx.WriteRune('}')
}

func defaultFieldMarshaler[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) unsafeFieldMarshaler[Ctx] {
	var m unsafeMarshal[Ctx]
	ctx.lazyGetter(t, &m)
	return func(name string, ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		return name, m(ptr, ctx)
	}
}
