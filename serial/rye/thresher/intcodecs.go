package thresher

import (
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
	"golang.org/x/exp/constraints"
)

func intCodec[I constraints.Signed]() {
	t := reflector.Type[I]()
	codecs[t] = &codec{
		enc: func(v any, s compact.Serializer) {
			s.CompactInt64(int64(v.(I)))
		},
		size: func(v any) uint64 {
			return compact.SizeInt64(int64(v.(I)))
		},
		encodingID: intEncID,
	}
	decoders[typeEncoding{
		encID: string(intEncID),
		t:     t,
	}] = func(d compact.Deserializer) any {
		return I(d.CompactInt64())
	}
}

func uintCodec[U constraints.Unsigned]() {
	t := reflector.Type[U]()
	codecs[t] = &codec{
		enc: func(v any, s compact.Serializer) {
			s.CompactUint64(uint64(v.(U)))
		},
		size: func(v any) uint64 {
			return compact.SizeUint64(uint64(v.(U)))
		},
		encodingID: uintEncID,
	}
	decoders[typeEncoding{
		encID: string(uintEncID),
		t:     t,
	}] = func(d compact.Deserializer) any {
		return U(d.CompactUint64())
	}
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
