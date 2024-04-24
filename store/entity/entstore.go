package entity

import (
	"bytes"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
)

const (
	// ErrKeyNotFound is returned by Load when the key is not found the Store.
	ErrKeyNotFound = lerr.Str("key not found")
	// ErrIndexNotFound
	ErrIndexNotFound       = lerr.Str("index not found")
	ErrIndexNameBlank      = lerr.Str("index name was blank")
	ErrIndexAlreadyDefined = lerr.Str("an index with the given name already exists")
)

// EntStore provides methods for saving and retreiving and Entity from a Store.
// The Serializer and Deserializer will be used for the Entity Value while
// EntKey will be used for the key.
type EntStore[T Entity] struct {
	Init func() T
	store.Store
	serial.Serializer
	serial.Deserializer
	IdxStore store.Store
	indexes  lmap.Map[string, Indexer[T]]
}

// Load an entity. This requires that EntKey returns the key of the value to be
// loaded.
func (es *EntStore[T]) Load(ent T) error {
	r := es.Store.Get(ent.EntKey())
	if r.Value == nil {
		return ErrKeyNotFound
	}

	return es.Deserialize(ent, r.Value)
}

func (es *EntStore[T]) AddIndex(name string, multi bool, fn func(T) []byte) StoreIndex[T] {
	return es.AddIndexer(BaseIndexer[T]{
		IndexName: name,
		Fn:        fn,
		M:         multi,
	})
}

func (es *EntStore[T]) AddIndexer(idx Indexer[T]) StoreIndex[T] {
	if es.indexes == nil {
		es.indexes = make(lmap.Map[string, Indexer[T]])
	}
	name := idx.Name()
	if name == "" {
		lerr.Panic(ErrIndexNameBlank)
	}
	if _, found := es.indexes[name]; found {
		lerr.Panic(ErrIndexAlreadyDefined)
	}
	es.indexes[name] = idx
	return StoreIndex[T]{
		idx: idx,
		es:  es,
	}
}

// Get an entity by key.
func (es *EntStore[T]) Get(key []byte) (found bool, ent T, err error) {
	r := es.Store.Get(key)
	found = r.Value != nil
	if !found {
		return
	}

	ent = es.Init()
	err = es.Deserialize(ent, r.Value)
	if err != nil {
		return
	}
	if set, ok := Entity(ent).(EntKeySetter); ok {
		set.SetEntKey(key)
	}

	return
}

// GetSlice returns all entities in the Store if fn is nil and if fn is defined
// it returns all entities for which fn returns true given their key. If buf
// is provided, it will be used as the return slice.
func (es *EntStore[T]) GetSlice(fn filter.Filter[[]byte], buf []T) ([]T, error) {
	i := store.NewIter(es.Store)
	i.Filter = fn
	return es.GetIter(i, buf)
}

func (es *EntStore[T]) GetIter(i liter.Iter[[]byte], buf []T) ([]T, error) {
	var m lerr.Many
	return slice.Transform(i, func(id []byte, _ int) (T, bool) {
		found, e, err := es.Get(id)
		m = m.Add(err)
		return e, found
	}), m.Cast()
}

// Put writes an entity to the store. If buf is provided, it will be used for
// Serialization.
func (es *EntStore[T]) Put(ent T, buf []byte) ([]byte, error) {
	v, err := es.Serialize(ent, buf)
	if err != nil {
		return nil, err
	}
	ek := ent.EntKey()
	err = es.Store.Put(ek, v)
	if err != nil {
		return nil, err
	}

	if es.IdxStore != nil {
		es.updateIndexes(ek, ent)
	}

	return v, nil
}

func (es *EntStore[T]) updateIndexes(ek []byte, ent T) {
	keysBkt := lerr.Must(es.IdxStore.Store([]byte("_keys")))
	prevKeys := make(map[string][]byte)
	if r := keysBkt.Get(ek); r.Found {
		es.Deserializer.Deserialize(&prevKeys, r.Value)
	}

	keys := make(map[string][]byte, len(es.indexes))
	for n, i := range es.indexes {
		k := i.IndexKey(ent)
		pk, ok := prevKeys[n]
		if !ok || !bytes.Equal(k, pk) {
			keys[n] = k
		}
	}

	for n, k := range keys {
		es.updateIdx(ek, k, prevKeys[n], es.indexes[n])
	}
	ks := lerr.Must(es.Serializer.Serialize(keys, nil))
	keysBkt.Put(ek, ks)
}

func (es *EntStore[T]) Index(name string) (StoreIndex[T], bool) {
	idx, found := es.indexes[name]
	si := StoreIndex[T]{
		es:  es,
		idx: idx,
	}
	return si, found
}

type StoreIndex[E Entity] struct {
	idx Indexer[E]
	es  *EntStore[E]
}

func (si StoreIndex[E]) Lookup(idxKey []byte) (entIds liter.Iter[[]byte], err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	entIds = si.es.idx(idxKey, si.idx)
	return
}

func (si StoreIndex[E]) LookupEnts(idxKey []byte) (ents slice.Slice[E], err error) {
	entIds, err := si.Lookup(idxKey)
	if err == nil {
		ents, err = si.es.GetIter(entIds, nil)
	}
	return
}

func (si StoreIndex[E]) Search(f filter.Filter[[]byte]) (idxKeys liter.Iter[[]byte], err error) {
	i := store.NewIter(si.es.getIdxBkt(si.idx))
	i.Filter = f
	return i, nil
}

func (si StoreIndex[E]) SearchEnts(f filter.Filter[[]byte]) (ents slice.Slice[E], err error) {
	s, _ := si.Search(f)
	return si.MultiEntLookup(s)
}

func (si StoreIndex[E]) MultiLookup(idxKeys liter.Iter[[]byte]) (entIds liter.Iter[[]byte], err error) {
	return &multiKeyLookup[E]{
		keys: idxKeys,
		si:   si,
	}, nil
}

func (si StoreIndex[E]) MultiEntLookup(idxKeys liter.Iter[[]byte]) (ents slice.Slice[E], err error) {
	entIds, err := si.MultiLookup(idxKeys)
	if err == nil {
		ents, err = si.es.GetIter(entIds, nil)
	}
	return
}

type multiKeyLookup[E Entity] struct {
	keys   liter.Iter[[]byte]
	entIds liter.Iter[[]byte]
	si     StoreIndex[E]
	i      int
}

func (mkl *multiKeyLookup[E]) Next() (entID []byte, done bool) {
	if mkl.entIds == nil || mkl.entIds.Done() {
		k, keysDone := mkl.keys.Next()
		if keysDone {
			return nil, true
		}
		mkl.entIds, _ = mkl.si.Lookup(k)
		entID, done = mkl.entIds.Cur()
		if done {
			return mkl.Next()
		}
		mkl.i++
		return
	}
	mkl.i++
	entID, _ = mkl.entIds.Next()
	return entID, false
}

func (mkl *multiKeyLookup[E]) Cur() (entID []byte, done bool) {
	if mkl.entIds == nil {
		k, keysDone := mkl.keys.Cur()
		for {
			if keysDone {
				return nil, true
			}
			mkl.entIds, _ = mkl.si.Lookup(k)
			entID, done = mkl.entIds.Cur()
			if !done {
				return
			}
			k, keysDone = mkl.keys.Next()
		}
	}
	entID, _ = mkl.entIds.Cur()
	return entID, mkl.keys.Done()
}
func (mkl *multiKeyLookup[E]) Done() bool {
	return mkl.keys.Done()
}
func (mkl *multiKeyLookup[E]) Idx() int {
	return mkl.i
}
