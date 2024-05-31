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
