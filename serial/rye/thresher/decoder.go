package thresher

import (
	"bytes"
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector/ltype"
)

type decoder func(d compact.Deserializer) any

type typeEncoding struct {
	encID string
	t     reflect.Type
}

var decoders = lmap.Map[typeEncoding, decoder]{
	{
		encID: string(compactSliceEncID),
		t:     ltype.String,
	}: func(d compact.Deserializer) any {
		return d.CompactString()
	},
	{
		encID: string(byteEncID),
		t:     ltype.Bool,
	}: func(d compact.Deserializer) any {
		return d.Byte() == 1
	},
	{
		encID: string(compactSliceEncID),
	}: func(d compact.Deserializer) any {
		return d.CompactSlice()
	},
	{
		encID: string(byteEncID),
	}: func(d compact.Deserializer) any {
		return d.Byte()
	},
}

// getDecdoer and getBaseDecoder are confusing me
// both Struct and slice should either always encode thier ID
// or only encode when base...

func getDecoder(t reflect.Type, id []byte) decoder {
	te := typeEncoding{
		encID: string(id),
		t:     t,
	}
	dec, found := decoders[te]
	if found {
		return dec
	}
	str := t.String()
	_ = str

	encoding := id
	var d compact.Deserializer
	if len(encoding) > 2 {
		d = compact.NewDeserializer(store[string(encoding)])
		encoding = d.CompactSlice()
	}

	if bytes.Equal(structEncID, encoding) {
		if k := t.Kind(); !(k == reflect.Invalid || k == reflect.Struct) {
			panic("should be struct")
		}
		dec = makeStructDecoder(t, d)
	} else if bytes.Equal(sliceEncID, encoding) {
		if k := t.Kind(); !(k == reflect.Invalid || k == reflect.Slice) {
			panic("should be slice")
		}
		dec = makeSliceDecoder(t, d)
	} else if bytes.Equal(compactSliceEncID, encoding) {
		if k := t.Kind(); !(k == reflect.Invalid || k == reflect.Pointer) {
			panic("should be pointer")
		}
		dec = pointerDecoder(t)
	} else {
		panic("deocder not found")
	}
	decoders[te] = dec
	return dec
}

func addDecoder(t reflect.Type, id []byte, d decoder) {
	decoders[typeEncoding{
		encID: string(id),
		t:     t,
	}] = d
}
