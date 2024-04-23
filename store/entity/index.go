package entity

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/liter"
)

type Indexer[E Entity] interface {
	Name() string
	IndexKey(E) []byte
	Multi() bool
}

type BaseIndexer[E Entity] struct {
	IndexName string
	Fn        func(E) []byte
	M         bool
}

func (bi BaseIndexer[E]) Name() string {
	return bi.IndexName
}
func (bi BaseIndexer[E]) IndexKey(e E) []byte {
	return bi.Fn(e)
}
func (bi BaseIndexer[E]) Multi() bool {
	return bi.M
}

var oneByte = []byte{0}

func (es EntStore[E]) updateIdx(ek, ik, pk []byte, idx Indexer[E]) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	bkt := lerr.Must(es.IdxStore.Store([]byte(idx.Name())))
	if idx.Multi() {
		es.deleteMultiKey(bkt, ek, pk)
		lerr.Must(bkt.Store(ik)).Put(ek, oneByte)
	} else {
		if pk != nil {
			bkt.Delete(pk)
		}
		bkt.Put(ik, ek)
	}
	return
}

func (es EntStore[E]) deleteMultiKey(bkt store.Store, ek, pk []byte) {
	if pk == nil {
		return
	}
	pkBkt := lerr.Must(bkt.Store(pk))
	lerr.Panic(pkBkt.Delete(ek))
	if pkBkt.Next(nil) == nil {
		lerr.Panic(bkt.Delete(pk))
	}
}

func (es EntStore[E]) idx(k []byte, idx Indexer[E]) (ids liter.Iter[[]byte]) {
	bkt := lerr.Must(es.IdxStore.Store([]byte(idx.Name())))
	r := bkt.Get(k)
	if !r.Found {
		panic(ErrKeyNotFound)
	}
	if idx.Multi() {
		ids = store.NewIter(r.Store)
	} else {
		ids = slice.New([][]byte{r.Value}).Iter()
	}

	return
}
