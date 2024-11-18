package ljson

import (
	"fmt"
	"reflect"

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

type valFieldMarshaler[Ctx any] func(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode)

type fieldMarshal[Ctx any] struct {
	idx  int
	name string
	//t    reflect.Type
	*valFieldMarshaler[Ctx]
}

type structMarshaler[Ctx any] []fieldMarshal[Ctx]

func (tctx *TypesContext[Ctx]) fieldMarshal(t reflect.Type, f reflect.StructField, self *valFieldMarshaler[Ctx]) {
	key := FieldKey{
		Type: t,
		Name: f.Name,
	}
	fm, found := tctx.fieldMarshalers[key]
	if found {
		*self = fm
		return
	}
	*self = func(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		fm, found := tctx.fieldMarshalers[key]
		if !found {
			fm = defaultFieldMarshaler(f.Type, tctx)
			tctx.fieldMarshalers[key] = fm
		}

		*self = fm
		return (*self)(name, v, ctx)
	}
}

func (tctx *TypesContext[Ctx]) buildStructMarshal(t reflect.Type) structMarshaler[Ctx] {
	return tctx.buildStructMarshalRecurse(t, 0, t.NumField(), 0)
}

func (tctx *TypesContext[Ctx]) buildStructMarshalRecurse(t reflect.Type, i, n, ln int) structMarshaler[Ctx] {
	if i == n {
		return make(structMarshaler[Ctx], ln)
	}
	f := t.Field(i)
	if !f.IsExported() {
		return tctx.buildStructMarshalRecurse(t, i+1, n, ln)
	}
	var fm valFieldMarshaler[Ctx]
	tctx.fieldMarshal(t, f, &fm)
	if fm == nil {
		return tctx.buildStructMarshalRecurse(t, i+1, n, ln)
	}
	sm := tctx.buildStructMarshalRecurse(t, i+1, n, ln+1)
	sm[ln] = fieldMarshal[Ctx]{
		idx:               i,
		name:              f.Name,
		valFieldMarshaler: &fm,
	}
	return sm
}

func (sm structMarshaler[Ctx]) valMarshal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	out := make(StructWriter, 0, len(sm))
	for _, fm := range sm {
		var fw FieldWriter

		var name string
		name, fw.Value = ctx.guardFieldMarshaler(fm.name, v.Field(fm.idx), *fm.valFieldMarshaler)
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

func defaultFieldMarshaler[Ctx any](t reflect.Type, ctx *TypesContext[Ctx]) valFieldMarshaler[Ctx] {
	var m valMarshaler[Ctx]
	ctx.get(t, &m)
	return func(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		return name, m(v, ctx)
	}
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
		panic("types do not match")
	}
	ctx.fieldMarshalers[key] = func(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
		t := v.Interface().(T)
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
			tctx.fieldMarshalers[k] = nil
		}
	}
}

// OmitEmpty for specified fields.
func (tctx *TypesContext[Ctx]) OmitEmpty(structKeys StructKeys, fieldNames ...string) {
	for _, name := range fieldNames {
		key := structKeys[name]
		st, ok := key.Field()
		if !ok {
			panic(fmt.Errorf("could not find %s on type %s", key.Name, key.Type.String()))
		}

		var vm valMarshaler[Ctx]
		tctx.get(st.Type, &vm)

		tctx.fieldMarshalers[key] = func(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
			if v.IsZero() {
				return "", nil
			}
			return name, vm(v, ctx)
		}
	}
}
