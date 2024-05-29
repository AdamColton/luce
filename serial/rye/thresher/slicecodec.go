package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/serial/rye/compact"
)

var baseSliceCodec *codec

func initBaseSliceCodec() {
	baseSliceCodec = &codec{
		enc: func(i any, s compact.Serializer) {
			v := reflect.ValueOf(i)
			ln := v.Len()
			s.Uint64(uint64(v.Len()))
			c := getCodec(v.Type().Elem())
			for i := 0; i < ln; i++ {
				c.enc(v.Index(i).Interface(), s)
			}
		},
		size: func(i any) uint64 {
			v := reflect.ValueOf(i)
			c := getCodec(v.Type().Elem())
			ln := v.Len()
			var out uint64 = 8 // for size
			for i := 0; i < ln; i++ {
				out += c.size(v.Index(i).Interface())
			}
			return out
		},
		roots: func(v reflect.Value) []*rootObj {
			c := getCodec(v.Type().Elem())
			if c.roots == nil {
				return nil
			}
			ln := v.Len()
			var out []*rootObj
			for i := 0; i < ln; i++ {
				e := v.Index(i)
				rts := c.roots(e)
				out = append(out, rts...)
			}
			return out
		},
	}
}

var sliceCodecs = lmap.Map[reflect.Type, *codec]{}

func getSliceCodec(t reflect.Type) *codec {
	c, found := sliceCodecs[t]
	if !found {
		c = &codec{
			enc:   baseSliceCodec.enc,
			size:  baseSliceCodec.size,
			roots: baseSliceCodec.roots,
			dec: func(d compact.Deserializer) any {
				ln := int(d.Uint64())
				s := reflect.MakeSlice(t, ln, ln)
				c := getCodec(t.Elem())
				for i := 0; i < ln; i++ {
					s.Index(i).Set(reflect.ValueOf(c.dec(d)))
				}
				return s.Interface()
			},
		}
		sliceCodecs[t] = c
	}
	return c
}
