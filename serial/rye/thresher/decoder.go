package thresher

import (
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
}

func getDecoder(t reflect.Type, id []byte) decoder {
	te := typeEncoding{
		encID: string(id),
		t:     t,
	}
	d, found := decoders[te]
	if found {
		return d
	}
	str := t.String()
	_ = str
	switch t.Kind() {
	case reflect.Struct:
		d = makeStructDecoder(t, id)
	case reflect.Slice:
		d = makeSliceDecoder(t, id)
	case reflect.Pointer:
		d = pointerDecoder(t)
	default:
		panic("deocder not found")
	}
	decoders[te] = d
	return d
}

func addDecoder(t reflect.Type, id []byte, d decoder) {
	decoders[typeEncoding{
		encID: string(id),
		t:     t,
	}] = d
}
