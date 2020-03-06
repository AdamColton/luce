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

func (dl *dataLeaf) size() int { return len(dl.data) }

func (*dataLeaf) Count() int                     { return 1 }
func (*dataLeaf) Depth() int                     { return 0 }
func (dl *dataLeaf) maxIdx() uint32              { return dl.idx }
func (dl *dataLeaf) have(idxs []uint32) []uint32 { return append(idxs, dl.idx) }

func uint32ToSlice(u uint32) []byte {
	out := make([]byte, 4)
	for i := 0; u > 0; i++ {
		out[i] = byte(u)
		u >>= 8
	}
	return out
}
