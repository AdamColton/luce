package merkle

import (
	"hash"

	"github.com/adamcolton/luce/serial/rye"
)

type dataLeaf struct {
	data, digest []byte
	idx          uint32
}

func newDataLeaf(data []byte, idx uint32, h hash.Hash) *dataLeaf {
	h.Reset()
	buf := make([]byte, 4)
	rye.Serialize.Uint32(buf, idx)
	h.Write(buf)
	h.Write(data)

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

func (dl *dataLeaf) Leaves() int {
	return 1
}
