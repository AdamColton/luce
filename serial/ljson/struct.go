package ljson

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/reflector"
)

// FieldKey is used to identify a Field by name on a Type.
type FieldKey struct {
	reflect.Type
	Name string
}

// Field returns the StructField the FieldKey represents.
func (fk FieldKey) Field() (reflect.StructField, bool) {
	return fk.Type.FieldByName(fk.Name)
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
		name, fw.Value = ctx.guardFieldMarshaler(fm.t, fm.name, fp, *fm.unsafeFieldMarshaler)
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

// StructKeys holds FieldKeys for all exported fields.
type StructKeys map[string]FieldKey

// GetFieldKeys returns the StructKeys for a struct.
func GetFieldKeys[T any]() StructKeys {
	t := reflector.Type[T]()
	if t.Kind() != reflect.Struct {
		return nil
	}
	n := t.NumField()
	out := make(map[string]FieldKey, n)
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.IsExported() {
			out[f.Name] = FieldKey{
				Type: t,
				Name: f.Name,
			}
		}
	}
	return out
}

// FieldMarshal is used to marshal a field on a struct. If either the WriteNode
// is nil or the string is empty, the field will not be written to the json
// document.
type FieldMarshal[T, Ctx any] func(name string, v T, ctx *MarshalContext[Ctx]) (string, WriteNode, error)

// AddFieldMarshal for the given key.
func AddFieldMarshal[T, Ctx any](key FieldKey, fm FieldMarshal[T, Ctx], ctx *TypesContext[Ctx]) {
	sf, found := key.Field()
	if !found {
		panic("key does not exist")
	}
	if sf.Type != reflector.Type[T]() {
		panic("types do not matcch")
	}
	ctx.fieldMarshal[key] = func(name string, ptr unsafe.Pointer, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		t := *(*T)(ptr)
		if fm == nil {
			return "", nil
		}
		name, wn, err := fm(name, t, ctx)
		lerr.Panic(err)
		return name, wn
	}
}

// OmitFields from a struct.
func (tctx *TypesContext[Ctx]) OmitFields(structKeys StructKeys, fieldNames ...string) {
	for _, n := range fieldNames {
		k, found := structKeys[n]
		if found {
			tctx.fieldMarshal[k] = nil
		}
	}
}
