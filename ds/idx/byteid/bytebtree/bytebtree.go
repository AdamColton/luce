package bytebtree

import (
	"bytes"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/sliceidx"

	"github.com/google/btree"
)

type byteIdxBTree struct {
	b  *btree.BTree
	si sliceidx.SliceIdx
}

// New fulfills byteid.Factory. It returns an instance of byteid.Index that
// stores the mapping from id to index in a btree using github.com/google/btree.
func New(sliceLen int) byteid.Index {
	return &byteIdxBTree{
		b:  btree.New(3),
		si: sliceidx.New(sliceLen),
	}
}

// SliceLen of the Indexed slice.
func (b *byteIdxBTree) SliceLen() int {
	return b.si.SliceLen
}

// SetSliceLen can be used to grow the slice.
func (b *byteIdxBTree) SetSliceLen(newlen int) {
	b.si.SetSliceLen(newlen)
}

// Insert an ID. The first value returned is the index and the bool
// indicates if an append is required.
func (b *byteIdxBTree) Insert(id []byte) (int, bool) {
	if idx, found := b.Get(id); found {
		return idx, false
	}

	var app bool
	e := entry{
		id: id,
	}
	e.idx, app = b.si.NextIdx()

	b.b.ReplaceOrInsert(e)
	return e.idx, app
}

// Get by ID. If not found it should return (-1,false). If it is found the
// first value is the index and the second value is True.
func (b *byteIdxBTree) Get(id []byte) (int, bool) {
	idx, found := -1, false
	b.b.AscendGreaterOrEqual(wrap(id), func(i btree.Item) bool {
		ie := i.(entry)
		if bytes.Compare(id, ie.id) == 0 {
			found = true
			idx = ie.idx
		}
		return false
	})

	return idx, found
}

// Delete by ID. Removes the ID from the index, the value should be
// recycled. This should be called before removing the value from the slice.
func (b *byteIdxBTree) Delete(id []byte) (int, bool) {
	e := b.b.Delete(wrap(id))
	if e == nil {
		return -1, false
	}
	idx := e.(entry).idx
	b.si.Recycle(idx)
	return idx, true
}

// Next ID after the ID given
func (b *byteIdxBTree) Next(id []byte) ([]byte, int) {
	var after []byte
	afterIdx := -1
	b.b.AscendGreaterOrEqual(wrap(id), func(i btree.Item) bool {
		ie := i.(entry)
		if !bytes.Equal(ie.id, id) {
			after = ie.id
			afterIdx = ie.idx
			return false
		}
		return true
	})
	return after, afterIdx
}
