package thresher

import (
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
	"golang.org/x/exp/constraints"
)

func intCodec[I constraints.Signed]() {
	t := reflector.Type[I]()
	encoders[t] = &encoder{
		encode: func(v any, s compact.Serializer, base bool) {
			if base {
				s.CompactSlice(intEncID)
			}
			s.CompactInt64(int64(v.(I)))
		},
		size: func(v any, base bool) uint64 {
			size := compact.SizeInt64(int64(v.(I)))
			if base {
				size += intEncIDSize
			}
			return size
		},
		encodingID: intEncID,
	}
	dec := func(d compact.Deserializer) any {
		return I(d.CompactInt64())
	}
	addDecoder(t, intEncID, dec)
	addDecoder(nil, intEncID, dec)
}

func uintCodec[U constraints.Unsigned]() {
	t := reflector.Type[U]()
	encoders[t] = &encoder{
		encode: func(v any, s compact.Serializer, base bool) {
			if base {
				s.CompactSlice(uintEncID)
			}
			s.CompactUint64(uint64(v.(U)))
		},
		size: func(v any, base bool) uint64 {
			size := compact.SizeUint64(uint64(v.(U)))
			if base {
				size += uintEncIDSize
			}
			return size
		},
		encodingID: uintEncID,
	}
	dec := func(d compact.Deserializer) any {
		return U(d.CompactUint64())
	}
	addDecoder(t, uintEncID, dec)
	addDecoder(nil, uintEncID, dec)
}

func initIntCodecs() {
	intCodec[int]()
	intCodec[int16]()
	intCodec[int32]()
	intCodec[int64]()
	uintCodec[uint]()
	uintCodec[uint16]()
	uintCodec[uint32]()
	uintCodec[uint64]()
}
