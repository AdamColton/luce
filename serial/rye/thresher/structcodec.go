package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/serial/rye/compact"
)

type fieldCodec struct {
	idx int
	*codec
}

type structCodec struct {
	fields []fieldCodec
	reflect.Type
}

func (sc *structCodec) enc(i any, s compact.Serializer) {
	v := reflect.ValueOf(i)
	for _, fc := range sc.fields {
		f := v.Field(fc.idx).Interface()
		fc.enc(f, s)
	}
}

func (sc *structCodec) dec(d compact.Deserializer) any {
	srct := reflect.New(sc.Type).Elem()
	for _, fc := range sc.fields {
		idx := fc.idx
		i := fc.dec(d)
		if i != nil {
			fv := reflect.ValueOf(i)
			srct.Field(idx).Set(fv)
		}
	}
	return srct.Interface()
}

func (sc *structCodec) size(i any) (sum uint64) {
	v := reflect.ValueOf(i)
	for _, fc := range sc.fields {
		f := v.Field(fc.idx).Interface()
		sum += fc.size(f)
	}
	return
}

func (sc *structCodec) roots(v reflect.Value) (out []*rootObj) {
	for _, fc := range sc.fields {
		if fc.roots != nil {
			f := v.Field(fc.idx)
			out = append(out, fc.roots(f)...)
		}
	}
	return
}

func makeStructCodec(t reflect.Type) *codec {
	c := &codec{}
	codecs[t] = c
	sc := &structCodec{
		Type:   t,
		fields: fieldsRecur(0, t.NumField(), t, 0),
	}
	c.enc = sc.enc
	c.dec = sc.dec
	c.size = sc.size
	c.roots = sc.roots
	return c
}

func fieldsRecur(i int, ln int, t reflect.Type, fields int) []fieldCodec {
	for ; i < ln; i++ {
		f := t.Field(i)
		if f.IsExported() {
			if c := getCodec(f.Type); c != nil {
				fcs := fieldsRecur(i+1, ln, t, fields+1)
				fcs[fields].idx = i
				fcs[fields].codec = c
				return fcs
			}
		}
	}
	return make([]fieldCodec, fields)
}
