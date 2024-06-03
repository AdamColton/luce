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

func getSliceEncoder(t reflect.Type) *encoder {
	sc, found := sliceEncoders[t]
	if !found {
		et := t.Elem()
		ec := getEncoder(et)
		edec := getDecoder(et, ec.encodingID)

		size := sliceEncIDSize + compact.Size(ec.encodingID)
		encoding := compact.MakeSerializer(int(size))
		encoding.CompactSlice(sliceEncID)
		encoding.CompactSlice(ec.encodingID)
		h := sha256.New()
		h.Write(encoding.Data)
		encID := h.Sum(nil)
		encIDSize := compact.Size(encID)

		sc = &encoder{
			encode: func(i any, s compact.Serializer, base bool) {
				if base {
					s.CompactSlice(encID)
				}
				v := reflect.ValueOf(i)
				ln := v.Len()
				s.CompactUint64(uint64(ln))
				for i := 0; i < ln; i++ {
					ec.encode(v.Index(i).Interface(), s, false)
				}
			},
			size: func(i any, base bool) uint64 {
				v := reflect.ValueOf(i)
				c := getEncoder(v.Type().Elem())
				ln := v.Len()
				out := compact.SizeUint64(uint64(v.Len()))
				if base {
					out += encIDSize
				}
				for i := 0; i < ln; i++ {
					out += c.size(v.Index(i).Interface(), false)
				}
				return out
			},
			roots:      sliceRoots,
			encodingID: encID,
		}
		sliceEncoders[t] = sc

		store[string(encID)] = encoding.Data
		dec := func(d compact.Deserializer) any {
			ln := int(d.CompactUint64())
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

func makeSliceDecoder(t reflect.Type, d compact.Deserializer) decoder {
	nilType := t.Kind() == reflect.Invalid
	var et reflect.Type
	if !nilType {
		et = t.Elem()
	}
	edec := getDecoder(et, d.CompactSlice())

	if nilType {
		return func(d compact.Deserializer) any {
			ln := int(d.CompactUint64())
			for i := 0; i < ln; i++ {
				edec(d)
			}
			return nil
		}
	}

	return func(d compact.Deserializer) any {
		ln := int(d.CompactUint64())
		s := reflect.MakeSlice(t, ln, ln)
		for i := 0; i < ln; i++ {
			s.Index(i).Set(reflect.ValueOf(edec(d)))
		}
		return s.Interface()
	}
}
