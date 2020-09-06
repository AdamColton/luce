package bytetree

import (
	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/byteid/sliceidx"
)

type node struct {
	children [256]*node
	rest     []byte
	idx      int
}
type byteIdxByteTree struct {
	root *node
	si   sliceidx.SliceIdx
}

// New fulfills byteid.Factory.
func New(sliceLen int) byteid.Index {
	return &byteIdxByteTree{
		root: &node{
			idx: -1,
		},
		si: sliceidx.New(sliceLen),
	}
}

func (bt *byteIdxByteTree) Insert(id []byte) (int, bool) {
	sr := bt.seek(id)
	if sr.found {
		return sr.idx, false
	}
	idx, app := bt.si.NextIdx()
	sr.insert(id, idx)
	return idx, app
}

func (bt *byteIdxByteTree) Get(id []byte) (int, bool) {
	sr := bt.seek(id)
	return sr.idx, sr.found
}

func (bt *byteIdxByteTree) Delete(id []byte) (int, bool) {
	return 0, false
}
func (bt *byteIdxByteTree) SliceLen() int {
	return bt.si.SliceLen
}
func (bt *byteIdxByteTree) SetSliceLen(newLen int) {
	bt.si.SetSliceLen(newLen)
}
func (bt *byteIdxByteTree) Next(id []byte) ([]byte, int) {
	return nil, 0
}
