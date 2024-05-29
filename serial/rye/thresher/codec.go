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

func init() {
	initPointerCoded()
	initBaseSliceCodec()
}
