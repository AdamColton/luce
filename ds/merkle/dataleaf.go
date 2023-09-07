package merkle

import (
	"hash"
)

type dataLeaf struct {
	data, digest []byte
	idx          uint32
}

func newDataLeaf(data []byte, idx uint32, h hash.Hash) *dataLeaf {
	h.Reset()
	h.Write(data)
	h.Write(uint32ToSlice(idx))
	return &dataLeaf{
		data:   data,
		digest: h.Sum(nil),
		idx:    idx,
	}
}

func (dl *dataLeaf) Digest() []byte {
	return dl.digest
}

func (dl *dataLeaf) Data() []byte {
	return dl.data
}
