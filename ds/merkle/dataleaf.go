package merkle

import (
	"hash"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/serial/rye"
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
	buf := slice.NewBuffer(dl.digest).Empty(cmpr.Max(4, h.Size()))

	h.Reset()
	rye.Serialize.Uint32(buf[:4], dl.idx)
	h.Write(buf[:4])
	h.Write(dl.data)
	dl.digest = h.Sum(buf)
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
