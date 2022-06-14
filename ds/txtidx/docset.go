package txtidx

import "github.com/adamcolton/luce/ds/bintrie"

type DocSet interface {
	Add(DocIDer)
	Has(DocIDer) bool
	Intersect(DocSet) DocSet
	Merge(with DocSet)
	Union(with DocSet) DocSet
	Len() int
	IDs() []DocID
}

type DocIDer interface {
	ID() DocID
}

type DocID uint32

func (id DocID) ID() DocID {
	return id
}

type docSet struct {
	t bintrie.Trie
}

func newDocSet() *docSet {
	return &docSet{
		t: bintrie.New(),
	}
}

func (ds *docSet) Len() int {
	return ds.t.Size()
}

func (ds *docSet) Add(di DocIDer) {
	id32 := uint32(di.ID())
	ds.t.Insert(id32)
}

func (ds *docSet) Has(di DocIDer) bool {
	id32 := uint32(di.ID())
	return ds.t.Has(id32)
}

func (ds *docSet) Intersect(with DocSet) DocSet {
	return ds.intersect(with.(*docSet))
}

func (ds *docSet) intersect(with *docSet) *docSet {
	return &docSet{
		t: bintrie.And(ds.t, with.t),
	}
}

func (ds *docSet) IntersectMerge(with DocSet) {
	ds.intersectMerge(with.(*docSet))
}

func (ds *docSet) intersectMerge(with *docSet) {
	ds.t.Union(with.t)
}

func (ds *docSet) Union(with DocSet) DocSet {
	return ds.union(with.(*docSet))
}

func (ds *docSet) union(with *docSet) *docSet {
	return &docSet{
		t: bintrie.Or(ds.t, with.t),
	}
}

func (ds *docSet) Merge(with DocSet) {
	ds.merge(with.(*docSet))
}

func (ds *docSet) merge(with *docSet) {
	ds.t.InsertTrie(with.t)
}

func (ds *docSet) Copy() DocSet {
	return ds.copy()
}

func (ds *docSet) copy() *docSet {
	return &docSet{
		t: ds.t.Copy(),
	}
}

func (ds *docSet) Delete(di DocIDer) {
	id32 := uint32(di.ID())
	ds.t.Delete(id32)
}

func (ds *docSet) IDs() []DocID {
	all := ds.t.All()
	out := make([]DocID, len(all))
	for i, b := range all {
		out[i] = DocID(b.ReadUint32())
	}
	return out
}
