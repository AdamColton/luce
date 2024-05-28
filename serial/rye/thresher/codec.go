package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

type codec struct {
	enc   func(i any, s compact.Serializer)
	dec   func(d compact.Deserializer) any
	size  func(i any) uint64
	roots func(i reflect.Value) []*rootObj
}

var codecs = map[reflect.Type]*codec{
	reflector.Type[string](): {
		enc: func(v any, s compact.Serializer) {
			s.CompactString(v.(string))
		},
		dec: func(d compact.Deserializer) any {
			return d.CompactString()
		},
		size: func(v any) uint64 {
			return compact.SizeString(v.(string))
		},
	},
	reflector.Type[bool](): {
		enc: func(v any, s compact.Serializer) {
			bit := byte(0)
			bol := v.(bool)
			if bol {
				bit = 1
			}
			s.Byte(bit)
		},
		dec: func(d compact.Deserializer) any {
			return d.Byte() == 1
		},
		size: func(v any) uint64 {
			return 1
		},
	},
	reflector.Type[int](): {
		enc: func(v any, s compact.Serializer) {
			s.CompactInt64(int64(v.(int)))
		},
		dec: func(d compact.Deserializer) any {
			return int(d.CompactInt64())
		},
		size: func(v any) uint64 {
			return compact.SizeInt64(int64(v.(int)))
		},
	},
}

func getCodec(t reflect.Type) *codec {
	c, found := codecs[t]
	if found {
		return c
	}

	switch t.Kind() {
	case reflect.Struct:
		c = makeStructCodec(t)
		codecs[t] = c
	case reflect.Pointer:
		c = pointerCodec
		codecs[t] = c
	}

	return c
}

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

var pointerCodec *codec

func init() {
	pointerCodec = &codec{
		enc: func(i any, s compact.Serializer) {
			ro := rootObjByV(reflect.ValueOf(i))
			s.CompactSlice(ro.getID())
		},
		dec: func(d compact.Deserializer) any {
			ro := getStoreByID(d.CompactSlice())
			if ro == nil {
				return nil
			}
			return ro.v.Interface()
		},
		size: func(i any) uint64 {
			ro := rootObjByV(reflect.ValueOf(i))
			return compact.Size(ro.getID())
		},
		roots: func(v reflect.Value) []*rootObj {
			ro := rootObjByV(v)
			if ro != nil {
				return []*rootObj{ro}
			}
			return nil
		},
	}
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
