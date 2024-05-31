package thresher

import (
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"golang.org/x/exp/constraints"
)

func intCodec[I constraints.Signed]() *codec {
	return &codec{
		enc: func(v any, s compact.Serializer) {
			s.CompactInt64(int64(v.(I)))
		},
		dec: func(d compact.Deserializer) any {
			return I(d.CompactInt64())
		},
		size: func(v any) uint64 {
			return compact.SizeInt64(int64(v.(I)))
		},
		encodingID: intEncID,
	}
}

func uintCodec[U constraints.Unsigned]() *codec {
	return &codec{
		enc: func(v any, s compact.Serializer) {
			s.CompactUint64(uint64(v.(U)))
		},
		dec: func(d compact.Deserializer) any {
			return U(d.CompactUint64())
		},
		size: func(v any) uint64 {
			return compact.SizeUint64(uint64(v.(U)))
		},
		encodingID: uintEncID,
	}
}

func initIntCodecs() {
	codecs[ltype.Int] = intCodec[int]()
	codecs[ltype.Int16] = intCodec[int16]()
	codecs[ltype.Int32] = intCodec[int32]()
	codecs[ltype.Int64] = intCodec[int64]()
	codecs[ltype.Uint] = uintCodec[uint]()
	codecs[ltype.Uint16] = uintCodec[uint16]()
	codecs[ltype.Uint32] = uintCodec[uint32]()
	codecs[ltype.Uint64] = uintCodec[uint64]()
}
