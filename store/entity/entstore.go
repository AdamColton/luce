package entity

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
)

// ErrKeyNotFound is returned by Load when the key is not found the Store.
const ErrKeyNotFound = lerr.Str("key not found")

// EntStore provides methods for saving and retreiving and Entity from a Store.
// The Serializer and Deserializer will be used for the Entity Value while
// EntKey will be used for the key.
type EntStore[T Entity] struct {
	Init func() T
	store.Store
	serial.Serializer
	serial.Deserializer
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
func (es *EntStore[T]) GetSlice(fn KeyFilter, buf []T) ([]T, error) {
	buf = buf[:0]
	for key := es.Store.Next(nil); key != nil; key = es.Store.Next(key) {
		if fn != nil && !fn(key) {
			continue
		}
		found, e, err := es.Get(key)
		if err != nil {
			return nil, err
		}
		if found {
			buf = append(buf, e)
		}
	}
	return buf, nil
}

// Put writes an entity to the store. If buf is provided, it will be used for
// Serialization.
func (es *EntStore[T]) Put(ent T, buf []byte) ([]byte, error) {
	v, err := es.Serialize(ent, buf)
	if err != nil {
		return nil, err
	}
	return v, es.Store.Put(ent.EntKey(), v)
}
