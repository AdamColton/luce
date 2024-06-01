package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

var (
	intEncID          = []byte{0, 1}
	uintEncID         = []byte{0, 2}
	byteEncID         = []byte{0, 3}
	compactSliceEncID = []byte{0, 4}
)

type codec struct {
	enc        func(i any, s compact.Serializer)
	size       func(i any) uint64
	roots      func(i reflect.Value) []*rootObj
	encodingID []byte
}

var codecs = map[reflect.Type]*codec{
	reflector.Type[string](): {
		enc: func(v any, s compact.Serializer) {
			s.CompactString(v.(string))
		},
		size: func(v any) uint64 {
			return compact.SizeString(v.(string))
		},
		encodingID: compactSliceEncID,
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
		size: func(v any) uint64 {
			return 1
		},
		encodingID: byteEncID,
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
		decoders[typeEncoding{
			encID: string(c.encodingID),
			t:     t,
		}] = pointerDecoder(t)
	case reflect.Slice:
		c = pointerCodec
		codecs[t] = c
	}

	return c
}

func getBaseCodec(t reflect.Type) (c *codec) {
	switch t.Kind() {
	case reflect.Struct:
		return getCodec(t)
	case reflect.Pointer:
		panic("cannnot get base codec of pointer")
	case reflect.Slice:
		c = getSliceCodec(t)
	}
	return
}
