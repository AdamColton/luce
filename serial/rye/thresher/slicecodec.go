package thresher

import (
	"bytes"
	"crypto/sha256"
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/serial/rye/compact"
)

var baseSliceCodec *codec

func initBaseSliceCodec() {
	baseSliceCodec = &codec{
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

var sliceHashPrefix = []byte{1, 1}

func getSliceCodec(t reflect.Type) *codec {
	sc, found := sliceCodecs[t]
	ec := getCodec(t.Elem())
	if !found {
		h := sha256.New()
		h.Write(sliceHashPrefix)
		h.Write(ec.encodingID)
		eid := h.Sum(nil)
		eidSize := compact.Size(eid)
		sc = &codec{
			enc: func(i any, s compact.Serializer) {
				v := reflect.ValueOf(i)
				ln := v.Len()
				s.CompactSlice(eid)
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
				var out uint64 = eidSize + 8 // for size
				for i := 0; i < ln; i++ {
					out += c.size(v.Index(i).Interface())
				}
				return out
			},
			roots: baseSliceCodec.roots,
			dec: func(d compact.Deserializer) any {
				id := d.CompactSlice()
				if !bytes.Equal(eid, id) {
					panic("encodingID does not match")
				}
				ln := int(d.Uint64())
				s := reflect.MakeSlice(t, ln, ln)
				for i := 0; i < ln; i++ {
					s.Index(i).Set(reflect.ValueOf(ec.dec(d)))
				}
				return s.Interface()
			},
			encodingID: eid,
		}
		sliceCodecs[t] = sc
		encodings[string(eid)] = ec.encodingID
	}
	return sc
}
