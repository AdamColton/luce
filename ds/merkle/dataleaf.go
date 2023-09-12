package merkle

import (
	"hash"
)

type dataLeaf struct {
	data, digest []byte
	idx          uint32
}

func newDataLeaf(data []byte, idx uint32) *dataLeaf {
	return &dataLeaf{
		data: data,
		idx:  idx,
	}
}

func (dl *dataLeaf) update(h hash.Hash) (ln, leaves int) {
	h.Reset()
	h.Write(uint32ToSlice(dl.idx))
	h.Write(dl.data)
	dl.digest = h.Sum(dl.digest)
	return len(dl.data), 1
}

func (dl *dataLeaf) Digest() []byte {
	return dl.digest
}

func (dl *dataLeaf) Data() []byte {
	return dl.data
}

func (dl *dataLeaf) Leaves() int {
	return 1
}

func (dl *dataLeaf) Len() int {
	return len(dl.data)
}

func (dl *dataLeaf) stitchData(in []byte) int {
	ln := copy(in, dl.data)
	dl.data = in[:ln]
	return ln
}
