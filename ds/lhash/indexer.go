package lhash

import (
	"hash"
	"unsafe"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
)

const (
	ErrBadKeysValue = lerr.Str("indexLen value must be greater than 0")
	ErrBadBitsValue = lerr.Str("bits value must be greater than 0")
)

type Uint interface {
	~uint | ~uint16 | ~uint32 | ~uint64
}

type Index[U Uint] []U

type Indexer[T any, U Uint] interface {
	Bits() int
	Index(t T) (Index[U], error)
}

type HashIndexer[T any, U Uint] struct {
	IndexLen uint32
	BitLen   byte
	Factory
	Converter func(t T) []byte
}

func (hi HashIndexer[T, U]) Bits() int {
	return int(hi.BitLen)
}

func (hi HashIndexer[T, U]) Index(t T) (Index[U], error) {
	bs := hi.Converter(t)
	return HashIndex[U](bs, hi.IndexLen, hi.BitLen, hi.Factory(), nil)
}

func HashIndex[U Uint](value []byte, indexLen uint32, bits byte, h hash.Hash, buf Index[U]) (Index[U], error) {
	if indexLen < 1 {
		return nil, ErrBadKeysValue
	}
	if bits < 1 {
		return nil, ErrBadBitsValue
	}

	b := slice.Buffer[U](buf)
	out := Index[U](b.Empty(int(indexLen))[:indexLen])
	modBits := uint64(1) << uint64(bits-1)

	// I need to efficiently consume the bits from the hash
	//

	for i := range out {
		h.Reset()
		_, err := h.Write(value)
		if err != nil {
			return nil, err
		}
		value = h.Sum(value[:0])
		u := *(*U)(unsafe.Pointer(&(value[0])))
		out[i] = u % U(modBits)
	}

	return out, nil
}
