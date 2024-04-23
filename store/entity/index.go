package entity

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
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

func (es EntStore[E]) updateIdx(ek, ik, pk []byte, idx Indexer[E]) error {
	bkt := lerr.Must(es.IdxStore.Store([]byte(idx.Name())))
	if idx.Multi() {
		if pk != nil {
			lerr.Must(bkt.Store(pk)).Delete(ek)
		}
		lerr.Must(bkt.Store(ik)).Put(ek, oneByte)
	} else {
		if pk != nil {
			bkt.Delete(pk)
		}
		bkt.Put(ik, ek)
	}
	return nil
}

func (es EntStore[E]) idx(k []byte, idx Indexer[E]) (ids slice.Slice[[]byte], err error) {
	var bkt store.Store
	bkt, err = es.IdxStore.Store([]byte(idx.Name()))
	if err != nil {
		return
	}
	r := bkt.Get(k)
	if !r.Found {
		err = ErrKeyNotFound
		return
	}
	if idx.Multi() {
		ids = store.Slice(r.Store)
	} else {
		ids = slice.New([][]byte{r.Value})
	}

	return
}
