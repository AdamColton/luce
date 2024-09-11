package scratchidx

import (
	"github.com/adamcolton/luce/ds/graph/rbtree"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye/compact"
	lgob "github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store/entity"
	"github.com/adamcolton/luce/util/filter"
)

type Index[K any, E entity.Entity] struct {
	tree  *rbtree.Tree[K, *lset.Set[string]]
	keyFn func(E) K
	prev  map[string]K
}

func NewIndex[K any, E entity.Entity](getKey func(E) K, cmpr filter.Compare[K], root []byte) *Index[K, E] {
	ptr := &entity.Reference[*rbtree.Node[K, *lset.Set[string]]]{
		ID: root,
	}
	return &Index[K, E]{
		tree:  rbtree.New(ptr, cmpr),
		keyFn: getKey,
		prev:  make(map[string]K),
	}
}

func (idx Index[K, E]) Store(b entity.Builder) (*entity.EntStore[*rbtree.Node[K, *lset.Set[string]]], error) {
	return entity.NewStore[*rbtree.Node[K, *lset.Set[string]]](b, "core", nil)
}

func (idx Index[K, E]) add(k K, id string) {
	n, found := idx.tree.Seek(k)
	idx.prev[id] = k
	if found {
		n.V.Add(id)
	} else {
		idx.tree.Add(k, lset.New(id))
	}
}

func (idx Index[K, E]) Remove(e E) {
	k := idx.keyFn(e)
	id := string(e.EntKey())
	idx.remove(k, id)
}

func (idx Index[K, E]) remove(k K, id string) {
	n, found := idx.tree.Seek(k)
	if found {
		n.V.Remove(id)
	}
	delete(idx.prev, id)
}

func (idx Index[K, E]) Update(e E) {
	id := string(e.EntKey())
	pk, rm := idx.prev[id]
	k := idx.keyFn(e)
	if rm {
		idx.remove(pk, id)
	}
	idx.add(k, id)
}

func (idx Index[K, E]) Get(k K) [][]byte {
	s, found := idx.tree.Seek(k)
	if !found {
		return nil
	}
	out := make([][]byte, 0, s.Len())
	s.V.Each(func(s string) (done bool) {
		out = append(out, []byte(s))
		return false
	})
	return out
}

func (idx *Index[K, E]) GobDecode(data []byte) (err error) {
	defer lerr.Recover(func(rerr error) {
		err = rerr
	})

	d := compact.NewDeserializer(data)
	lgob.Dec(d.CompactSlice(), idx.tree)
	lgob.Dec(d.CompactSlice(), &(idx.prev))
	return
}

func (idx *Index[K, E]) GobEncode() (b []byte, err error) {
	defer lerr.Recover(func(rerr error) {
		err = rerr
	})

	tree := lgob.Enc(idx.tree)
	ln := compact.Size(tree)

	prev := lgob.Enc(idx.prev)
	ln += compact.Size(prev)

	s := compact.MakeSerializer(int(ln))
	s.CompactSlice(tree)
	s.CompactSlice(prev)
	b = s.Data
	return
}
