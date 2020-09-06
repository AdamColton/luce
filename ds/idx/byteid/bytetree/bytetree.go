package bytetree

import (
	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/byteid/sliceidx"
)

type node struct {
	children   [256]*node
	childCount byte
	rest       []byte
	idx        int
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
	sr := bt.seek(id, false)
	if sr.found {
		return sr.idx, false
	}
	idx, app := bt.si.NextIdx()
	sr.insert(id, idx)
	return idx, app
}

func (bt *byteIdxByteTree) Get(id []byte) (int, bool) {
	sr := bt.seek(id, false)
	return sr.idx, sr.found
}

func (bt *byteIdxByteTree) Delete(id []byte) (int, bool) {
	sr := bt.seek(id, true)
	if !sr.found {
		return -1, false
	}
	idx := sr.idx
	sr.del(id)
	bt.si.Recycle(idx)
	return idx, true
}

func (bt *byteIdxByteTree) SliceLen() int {
	return bt.si.SliceLen
}
func (bt *byteIdxByteTree) SetSliceLen(newLen int) {
	bt.si.SetSliceLen(newLen)
}
func (bt *byteIdxByteTree) Next(id []byte) ([]byte, int) {
	sr := bt.seek(id, true)
	if id != nil && !sr.rightThenUp(sr.idInt(id)) {
		return nil, sr.idx
	}
	sr.downAndLeft()
	if sr.idx == -1 {
		return nil, sr.idx
	}
	return sr.value(), sr.idx
}
