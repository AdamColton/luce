package txtidx

import (
	"github.com/adamcolton/luce/ds/lset"
)

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
	t *lset.Set[DocID]
}

func newDocSet() *docSet {
	return &docSet{
		t: lset.New[DocID](),
	}
}

func (ds *docSet) Len() int {
	return ds.t.Len()
}

func (ds *docSet) Add(di DocIDer) {
	ds.t.Add(di.ID())
}

func (ds *docSet) Has(di DocIDer) bool {
	return ds.t.Contains(di.ID())
}

func (ds *docSet) Intersect(with DocSet) DocSet {
	return ds.intersect(with.(*docSet))
}

func (ds *docSet) intersect(with *docSet) *docSet {
	return &docSet{
		t: lset.And(ds.t, with.t),
	}
}

func (ds *docSet) IntersectMerge(with DocSet) {
	ds.intersectMerge(with.(*docSet))
}

func (ds *docSet) intersectMerge(with *docSet) {
	ds.t = lset.And(ds.t, with.t)
}

func (ds *docSet) Union(with DocSet) DocSet {
	return ds.union(with.(*docSet))
}

func (ds *docSet) union(with *docSet) *docSet {
	return &docSet{
		t: lset.Or(ds.t, with.t),
	}
}

func (ds *docSet) Merge(with DocSet) {
	ds.merge(with.(*docSet))
}

func (ds *docSet) merge(with *docSet) {
	ds.t.AddAll(with.t)
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
	ds.t.Remove(di.ID())
}

func (ds *docSet) IDs() []DocID {
	return ds.t.Slice()
}
