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

func (fk FieldKey) nameType() (string, reflect.Type) {
	sf, ok := fk.Field()
	if !ok {
		return "", nil
	}
	return sf.Name, sf.Type
}

type valFieldMarshaler[Ctx any] interface {
	marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode)
	nameType() (string, reflect.Type)
}

type fieldMarshal[Ctx any] struct {
	idx  int
	name string
	//t    reflect.Type
	valFieldMarshaler[Ctx]
}

type structMarshaler[Ctx any] []fieldMarshal[Ctx]

type deferGetFieldMarshal[Ctx any] struct {
	f    reflect.StructField
	t    reflect.Type
	self *valFieldMarshaler[Ctx]
}

func (dg deferGetFieldMarshal[Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	tctx := ctx.TypesContext
	key := FieldKey{
		Type: dg.t,
		Name: dg.f.Name,
	}
	fm, found := tctx.fieldMarshalers[key]
	if !found {
		fm = defaultFieldMarshaler(name, dg.f.Type, tctx)
		tctx.fieldMarshalers[key] = fm
	}

	*(dg.self) = fm
	return (*dg.self).marshalField(name, v, ctx)
}

func (dg deferGetFieldMarshal[Ctx]) nameType() (string, reflect.Type) {
	return dg.f.Name, dg.f.Type
}

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
	*self = deferGetFieldMarshal[Ctx]{
		f:    f,
		t:    t,
		self: self,
	}

}

func (tctx *TypesContext[Ctx]) buildStructMarshal(t reflect.Type) structMarshaler[Ctx] {
	return tctx.buildStructMarshalRecurse(t, 0, t.NumField(), 0)
}

func (tctx *TypesContext[Ctx]) buildStructMarshalRecurse(t reflect.Type, i, n, ln int) structMarshaler[Ctx] {
	if i == n {
		return tctx.attachFieldGenerators(t, ln)
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
		valFieldMarshaler: fm,
	}
	return sm
}

func (tctx *TypesContext[Ctx]) attachFieldGenerators(t reflect.Type, ln int) structMarshaler[Ctx] {
	fgs := tctx.fieldGenerators[t]
	out := make(structMarshaler[Ctx], ln+len(fgs))
	for i, fg := range fgs {
		fm := fg
		out[i+ln] = fieldMarshal[Ctx]{
			valFieldMarshaler: fm,
			idx:               -1,
		}
	}
	return out
}

func (sm structMarshaler[Ctx]) marshalVal(v reflect.Value, ctx *MarshalContext[Ctx]) WriteNode {
	out := make(StructWriter, 0, len(sm))
	for _, fm := range sm {
		var fw FieldWriter

		var name string
		if fm.idx >= 0 {
			name, fw.Value = ctx.guardFieldMarshaler(fm.name, v.Field(fm.idx), fm.valFieldMarshaler)
		} else {
			name, fw.Value = fm.valFieldMarshaler.marshalField("", v, ctx)
		}
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

type marshalValToField[Ctx any] struct {
	valMarshaler[Ctx]
	name string
	t    reflect.Type
}

func (vtf marshalValToField[Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	return name, vtf.marshalVal(v, ctx)
}

func (vtf marshalValToField[Ctx]) nameType() (string, reflect.Type) {
	return vtf.name, vtf.t
}

func defaultFieldMarshaler[Ctx any](name string, t reflect.Type, ctx *TypesContext[Ctx]) (out marshalValToField[Ctx]) {
	ctx.get(t, &(out.valMarshaler))
	return
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

func (fm FieldMarshal[T, Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	t := v.Interface().(T)
	if fm == nil {
		return "", nil
	}
	name, wn, err := fm(name, t, ctx)
	lerr.Panic(err)
	return name, wn
}

type fieldMarshalWrapper[T, Ctx any] struct {
	FieldMarshal[T, Ctx]
	FieldKey
}

// AddFieldMarshal for the given key.
func AddFieldMarshal[T, Ctx any](key FieldKey, fm FieldMarshal[T, Ctx], ctx *TypesContext[Ctx]) {
	sf, found := key.Field()
	if !found {
		panic("key does not exist")
	}
	if sf.Type != reflector.Type[T]() {
		panic("types do not match")
	}
	ctx.fieldMarshalers[key] = fieldMarshalWrapper[T, Ctx]{
		FieldKey:     key,
		FieldMarshal: fm,
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

type omitEmpty[Ctx any] struct {
	vm valMarshaler[Ctx]
	FieldKey
}

func (oe omitEmpty[Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	if v.IsZero() {
		return "", nil
	}
	return name, oe.vm.marshalVal(v, ctx)
}

// OmitEmpty for specified fields.
func (tctx *TypesContext[Ctx]) OmitEmpty(structKeys StructKeys, fieldNames ...string) {
	for _, name := range fieldNames {
		key := structKeys[name]
		st, ok := key.Field()
		if !ok {
			panic(fmt.Errorf("could not find %s on type %s", key.Name, key.Type.String()))
		}

		oe := omitEmpty[Ctx]{
			FieldKey: key,
		}
		tctx.get(st.Type, &(oe.vm))

		tctx.fieldMarshalers[key] = oe
	}
}

// FieldGenerator adds a field to On when marshaling.
type FieldGenerator[On, T, Ctx any] func(on On, ctx *MarshalContext[Ctx]) T

type marshalFieldGenerator[On, T, Ctx any] struct {
	um   valMarshaler[Ctx]
	fg   FieldGenerator[On, T, Ctx]
	name string
}

func (mfg marshalFieldGenerator[On, T, Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	ot := reflector.Type[On]()
	if ot.Kind() == reflect.Pointer {
		v = reflector.EnsurePointer(v)
	}
	on := v.Interface().(On)
	t := mfg.fg(on, ctx)
	return mfg.name, mfg.um.marshalVal(reflect.ValueOf(t), ctx)
}

func (mfg marshalFieldGenerator[On, T, Ctx]) nameType() (string, reflect.Type) {
	return mfg.name, reflector.Type[T]()
}

// GeneratedField adds the FieldGenerator to the TypesContext.
func GeneratedField[On, T, Ctx any](name string, fg FieldGenerator[On, T, Ctx], ctx *TypesContext[Ctx]) {
	mfg := marshalFieldGenerator[On, T, Ctx]{
		fg:   fg,
		name: name,
	}
	t := reflector.Type[T]()
	ctx.get(t, &(mfg.um))

	fgKey := reflector.Type[On]()
	if fgKey.Kind() == reflect.Pointer {
		fgKey = fgKey.Elem()
	}
	ctx.fieldGenerators[fgKey] = append(ctx.fieldGenerators[fgKey], mfg)
}

// ConditionalFunc uses a MarshalContext to generate a boolean.
type ConditionalFunc[Ctx any] func(ctx *MarshalContext[Ctx]) bool

type conditionalField[Ctx any] struct {
	cfn ConditionalFunc[Ctx]
	fm  valFieldMarshaler[Ctx]
	FieldKey
}

func (cf conditionalField[Ctx]) marshalField(name string, v reflect.Value, ctx *MarshalContext[Ctx]) (string, WriteNode) {
	if cf.cfn(ctx) {
		return cf.fm.marshalField(name, v, ctx)
	}
	return "", nil
}

// ConditionalFields omits the specified fields when the ConditionalFunc returns
// false.
func (tctx *TypesContext[Ctx]) ConditionalFields(cfn ConditionalFunc[Ctx], fieldKeys map[string]FieldKey, fieldNames ...string) {
	for _, n := range fieldNames {
		k, found := fieldKeys[n]
		if !found {
			continue
		}
		sf, found := k.Type.FieldByName(k.Name)
		if !found {
			continue
		}
		fm, found := tctx.fieldMarshalers[k]
		if !found {
			fm = defaultFieldMarshaler(n, sf.Type, tctx)
		}
		tctx.fieldMarshalers[k] = conditionalField[Ctx]{
			cfn:      cfn,
			fm:       fm,
			FieldKey: k,
		}
	}
}
