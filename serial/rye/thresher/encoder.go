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
	structEncID       = []byte{0, 5}

	// The fact these are being used is a code smell
	intEncIDSize          = compact.Size(intEncID)
	uintEncIDSize         = compact.Size(uintEncID)
	byteEncIDSize         = compact.Size(byteEncID)
	compactSliceEncIDSize = compact.Size(compactSliceEncID)
	structEncIDSize       = compact.Size(structEncID)
)

type encoder struct {
	encode     func(i any, s compact.Serializer, base bool)
	size       func(i any) uint64
	roots      func(i reflect.Value) []*rootObj
	encodingID []byte
}

var encoders = map[reflect.Type]*encoder{
	reflector.Type[string](): {
		encode: func(v any, s compact.Serializer, base bool) {
			if base {
				s.CompactSlice(compactSliceEncID)
			}
			s.CompactString(v.(string))
		},
		size: func(v any) uint64 {
			return compact.SizeString(v.(string))
		},
		encodingID: compactSliceEncID,
	},
	reflector.Type[bool](): {
		encode: func(v any, s compact.Serializer, base bool) {
			if base {
				s.CompactSlice(byteEncID)
			}
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

func getEncoder(t reflect.Type) *encoder {
	enc, found := encoders[t]
	if found {
		return enc
	}

	switch t.Kind() {
	case reflect.Struct:
		enc = getStructEncoder(t)
		encoders[t] = enc
	case reflect.Pointer:
		enc = pointerEncoder
		encoders[t] = enc
		addDecoder(t, enc.encodingID, pointerDecoder(t))
	case reflect.Slice:
		enc = pointerEncoder
		encoders[t] = enc
	}

	return enc
}

func getBaseEncoder(t reflect.Type) (enc *encoder) {
	switch t.Kind() {
	case reflect.Struct:
		return getEncoder(t)
	case reflect.Slice:
		enc = getSliceEncoder(t)
	default:
		panic("cannnot get base encoder of " + t.String())
	}
	return
}
