package hextree

import (
	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/byteid/sliceidx"
)

type hexTree struct {
	root        *low
	first, last *low
	si          sliceidx.SliceIdx
}

func New(slicelen int) byteid.Index {
	return &hexTree{
		root: &low{
			idx: -1,
		},
		si: sliceidx.New(slicelen),
	}
}

func (ht *hexTree) Insert(id []byte) (int, bool) {
	sr := ht.seek(id, false)
	if sr.found {
		return sr.l.idx, false
	}
	idx, app := ht.si.NextIdx()
	sr.insert(id, idx)
	return idx, app
}

func (ht *hexTree) Get(id []byte) (int, bool) {
	sr := ht.seek(id, false)
	return sr.l.idx, sr.found
}
func (ht *hexTree) Delete(id []byte) (int, bool) {
	sr := ht.seek(id, true)
	if !sr.found {
		return -1, false
	}
	idx := sr.l.idx
	sr.del(id)
	ht.si.Recycle(idx)
	return idx, true
}
func (ht *hexTree) SliceLen() int {
	return ht.si.SliceLen
}
func (ht *hexTree) SetSliceLen(newLen int) {
	ht.si.SetSliceLen(newLen)
}
func (ht *hexTree) Next(id []byte) ([]byte, int) {
	return nil, -1
}
