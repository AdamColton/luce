package bytebtree

import (
	"bytes"

	"github.com/adamcolton/luce/ds/idx/byteid"

	"github.com/google/btree"
)

type ByteBTree struct {
	b        *btree.BTree
	sliceLen int
	maxIdx   int
	recycle  []int
}

func New(sliceLen int) *ByteBTree {
	return &ByteBTree{
		b:        btree.New(3),
		sliceLen: sliceLen,
	}
}

func Factory(sliceLen int) byteid.Index {
	return byteid.Index(New(sliceLen))
}

func (b *ByteBTree) SliceLen() int {
	return b.sliceLen
}

func (b *ByteBTree) SetSliceLen(newlen int) {
	if newlen > b.sliceLen {
		b.sliceLen = newlen
	}
}

type entry struct {
	id  []byte
	idx int
}

func (e entry) Less(i btree.Item) bool {
	var id []byte
	switch i := i.(type) {
	case entry:
		id = i.id
	case wrap:
		id = i
	}
	return bytes.Compare(e.id, id) < 0
}

type wrap []byte

func (id wrap) Less(i btree.Item) bool {
	return bytes.Compare(id, i.(entry).id) == -1
}

func (b *ByteBTree) Insert(id []byte) (int, bool) {
	if idx, found := b.Get(id); found {
		return idx, false
	}

	e := entry{
		id: id,
	}
	app := false
	if ln := len(b.recycle); ln > 0 {
		e.idx = b.recycle[ln-1]
		b.recycle = b.recycle[:ln-1]
	} else {
		e.idx = b.maxIdx
		b.maxIdx++
		app = b.maxIdx > b.sliceLen
		if app {
			b.sliceLen = b.maxIdx
		}
	}
	b.b.ReplaceOrInsert(e)
	return e.idx, app
}

func (b *ByteBTree) Get(id []byte) (int, bool) {
	var found bool
	idx := -1
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

func (b *ByteBTree) Delete(id []byte) (int, bool) {
	e := b.b.Delete(wrap(id))
	if e == nil {
		return -1, false
	}
	idx := e.(entry).idx
	b.recycle = append(b.recycle, idx)
	return idx, true
}

func (b *ByteBTree) Next(id []byte) ([]byte, int) {
	var after []byte
	afterIdx := -1
	b.b.AscendGreaterOrEqual(wrap(id), func(i btree.Item) bool {
		ie := i.(entry)
		if bytes.Compare(ie.id, id) != 0 {
			after = ie.id
			afterIdx = ie.idx
			return false
		}
		return true
	})
	return after, afterIdx
}
