package thresher

import (
	"crypto/sha256"
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/serial/rye/compact"
)

func sliceRoots(v reflect.Value) []*rootObj {
	c := getEncoder(v.Type().Elem())
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
}

var sliceEncoders = lmap.Map[reflect.Type, *encoder]{}

var sliceHashPrefix = []byte{1, 1}

func getSliceEncoder(t reflect.Type) *encoder {
	sc, found := sliceEncoders[t]
	if !found {
		et := t.Elem()
		ec := getEncoder(et)
		edec := getDecoder(et, ec.encodingID)

		h := sha256.New()
		h.Write(sliceHashPrefix)
		h.Write(ec.encodingID)
		encID := h.Sum(nil)
		encIDSize := compact.Size(encID)
		sc = &encoder{
			encode: func(i any, s compact.Serializer, base bool) {
				if base {
					s.CompactSlice(encID)
				}
				v := reflect.ValueOf(i)
				ln := v.Len()
				s.Uint64(uint64(v.Len()))
				c := getEncoder(v.Type().Elem())
				for i := 0; i < ln; i++ {
					c.encode(v.Index(i).Interface(), s, false)
				}
			},
			size: func(i any) uint64 {
				v := reflect.ValueOf(i)
				c := getEncoder(v.Type().Elem())
				ln := v.Len()
				var out uint64 = encIDSize + 8 // for size
				for i := 0; i < ln; i++ {
					out += c.size(v.Index(i).Interface())
				}
				return out
			},
			roots:      sliceRoots,
			encodingID: encID,
		}
		sliceEncoders[t] = sc
		store[string(encID)] = ec.encodingID
		dec := func(d compact.Deserializer) any {
			ln := int(d.Uint64())
			s := reflect.MakeSlice(t, ln, ln)
			for i := 0; i < ln; i++ {
				s.Index(i).Set(reflect.ValueOf(edec(d)))
			}
			return s.Interface()
		}
		addDecoder(t, encID, dec)
	}
	return sc
}

func makeSliceDecoder(t reflect.Type, id []byte) decoder {
	edec := getDecoder(t.Elem(), store[string(id)])
	return func(d compact.Deserializer) any {
		ln := int(d.Uint64())
		s := reflect.MakeSlice(t, ln, ln)
		for i := 0; i < ln; i++ {
			s.Index(i).Set(reflect.ValueOf(edec(d)))
		}
		return s.Interface()
	}
}
