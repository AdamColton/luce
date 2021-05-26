package lusers

import (
	"bytes"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
)

// failFactory fails NestedStore for the one name in failOn, and otherwise
// delegates, optionally wrapping the real Sub-Store it returns for wrapOn.
type failFactory struct {
	store.NestedFactory
	failOn []byte
	wrapOn []byte
	wrap   func(store.NestedStore) store.NestedStore
}

func (f failFactory) NestedStore(name []byte) (store.NestedStore, error) {
	if f.failOn != nil && bytes.Equal(name, f.failOn) {
		return nil, lerr.Str("NestedStore failed")
	}
	s, err := f.NestedFactory.NestedStore(name)
	if err != nil || s == nil {
		return s, err
	}
	if f.wrapOn != nil && bytes.Equal(name, f.wrapOn) && f.wrap != nil {
		return f.wrap(s), nil
	}
	return s, nil
}

// errPutStore wraps a NestedStore and always fails Put. It implements every
// method explicitly: embedding store.NestedStore under its own type name
// would make the field itself (not a promoted method) answer to the
// selector "NestedStore", which is the method this type also needs to
// delegate.
type errPutStore struct{ inner store.NestedStore }

func (s errPutStore) Put(key, value []byte) error { return lerr.Str("Put failed") }
func (s errPutStore) Get(key []byte) store.Record { return s.inner.Get(key) }
func (s errPutStore) Next(key []byte) []byte      { return s.inner.Next(key) }
func (s errPutStore) Delete(key []byte) error     { return s.inner.Delete(key) }
func (s errPutStore) Len() int                    { return s.inner.Len() }
func (s errPutStore) NestedStore(name []byte) (store.NestedStore, error) {
	return s.inner.NestedStore(name)
}

// errNestStore wraps a NestedStore and always fails NestedStore (see
// errPutStore for why it can't just embed store.NestedStore).
type errNestStore struct{ inner store.NestedStore }

func (s errNestStore) Put(key, value []byte) error { return s.inner.Put(key, value) }
func (s errNestStore) Get(key []byte) store.Record { return s.inner.Get(key) }
func (s errNestStore) Next(key []byte) []byte      { return s.inner.Next(key) }
func (s errNestStore) Delete(key []byte) error     { return s.inner.Delete(key) }
func (s errNestStore) Len() int                    { return s.inner.Len() }
func (s errNestStore) NestedStore(name []byte) (store.NestedStore, error) {
	return nil, lerr.Str("NestedStore failed")
}
