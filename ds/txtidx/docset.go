package txtidx

type DocSet interface {
	Add(DocIDer)
	Has(DocIDer) bool
	Len() int
	Intersect(DocSet) DocSet
	Merge(with DocSet)
	Pop() (id DocID, done bool)
}

type DocIDer interface {
	ID() DocID
}

type DocID uint32

func (id DocID) ID() DocID {
	return id
}

type docSet struct {
	docs map[DocID]sig
}

func newDocSet() *docSet {
	return &docSet{
		docs: map[DocID]sig{},
	}
}

func (ds *docSet) Add(di DocIDer) {
	ds.docs[di.ID()] = sig{}
}

func (ds *docSet) Has(di DocIDer) bool {
	_, found := ds.docs[di.ID()]
	return found
}

func (ds *docSet) Len() int {
	return len(ds.docs)
}

func (ds *docSet) Intersect(with DocSet) DocSet {
	return ds.intersect(with.(*docSet))
}

func (ds *docSet) intersect(with *docSet) *docSet {
	out := map[DocID]sig{}
	iter, cmpr := ds.docs, with.docs
	if len(iter) > len(cmpr) {
		iter, cmpr = cmpr, iter
	}
	for di := range iter {
		_, found := cmpr[di]
		if found {
			out[di] = sig{}
		}
	}
	return &docSet{
		docs: out,
	}
}

func (ds *docSet) Merge(with DocSet) {
	ds.merge(with.(*docSet))
}

func (ds *docSet) merge(with *docSet) {
	for di := range with.docs {
		ds.docs[di] = sig{}
	}
}

func (ds *docSet) Copy() DocSet {
	return ds.copy()
}

func (ds *docSet) copy() *docSet {
	out := &docSet{
		docs: map[DocID]sig{},
	}
	out.merge(ds)
	return out
}

func (ds *docSet) Delete(di DocIDer) {
	delete(ds.docs, di.ID())
}

func (ds *docSet) Pop() (DocID, bool) {
	if len(ds.docs) == 0 {
		return DocID(MaxUint32), true
	}
	var out DocID
	for di := range ds.docs {
		out = di
		break
	}
	ds.Delete(out)
	return out, false
}
