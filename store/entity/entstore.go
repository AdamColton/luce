package entity

import (
	"bytes"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/iter"
)

const (
	// ErrKeyNotFound is returned by Load when the key is not found the Store.
	ErrKeyNotFound = lerr.Str("key not found")
	// ErrIndexNotFound
	ErrIndexNotFound = lerr.Str("index not found")
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
	Indexes  lmap.Map[string, Indexer[T]]
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

func (es *EntStore[T]) AddIndex(name string, multi bool, fn func(T) []byte) {
	es.AddIndexer(BaseIndexer[T]{
		IndexName: name,
		Fn:        fn,
		M:         multi,
	})
}

func (es *EntStore[T]) AddIndexer(idx Indexer[T]) {
	if es.Indexes == nil {
		es.Indexes = make(lmap.Map[string, Indexer[T]])
	}
	es.Indexes[idx.Name()] = idx
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

func (es *EntStore[T]) GetIter(i iter.Iter[[]byte], buf []T) ([]T, error) {
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

	keys := make(map[string][]byte, len(es.Indexes))
	for n, i := range es.Indexes {
		k := i.IndexKey(ent)
		pk, ok := prevKeys[n]
		if !ok || !bytes.Equal(k, pk) {
			keys[n] = k
		}
	}

	for n, k := range keys {
		es.updateIdx(ek, k, prevKeys[n], es.Indexes[n])
	}
	ks := lerr.Must(es.Serializer.Serialize(keys, nil))
	keysBkt.Put(ek, ks)
}

func (es *EntStore[T]) Index(name string, key []byte) (ents slice.Slice[T], err error) {
	defer func() {
		err, _ = recover().(error)
	}()
	idx, ok := es.Indexes[name]
	if !ok {
		panic(ErrIndexNotFound)
	}
	ids := es.idx(key, idx)
	ents, err = es.GetIter(ids, nil)

	return
}
