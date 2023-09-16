// Package ephemeral implements a store in memory. Nothing is persisted.
package ephemeral

import (
	"slices"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/store"
)

// estore is a store.FlatStore held in memory. The index maps each key to a
// position in records. It copies every key and value it is given or returns, so
// callers may reuse their buffers.
type estore struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	records    [][]byte
}

// Len is the number of keys in the store.
func (e *estore) Len() int {
	return e.idx.Len()
}

// Put stores a copy of value at key, replacing any existing value. It never
// returns an error.
func (e *estore) Put(key, value []byte) error {
	value = slices.Clone(value)
	idx, found := e.idx.Get(key)
	if found {
		e.records[idx] = value
		return nil
	}
	idx, app := e.idx.Insert(slices.Clone(key))
	if app {
		e.records = append(e.records, value)
	} else {
		e.records[idx] = value
	}
	return nil
}

// Get returns a copy of the value stored at key.
func (e *estore) Get(key []byte) store.Record {
	var out store.Record
	var idx int
	idx, out.Found = e.idx.Get(key)
	if out.Found {
		out.Value = slices.Clone(e.records[idx])
	}
	return out
}

// Next returns a copy of the smallest key greater than key, or of the lowest key
// if key is nil. It returns nil if there is no such key.
func (e *estore) Next(key []byte) []byte {
	next, _ := e.idx.Next(key)
	return slices.Clone(next)
}

// Delete removes key. Deleting a key that is not in the store is not an error.
func (e *estore) Delete(key []byte) error {
	idx, found := e.idx.Delete(key)
	if found {
		e.records[idx] = nil
	}
	return nil
}

// newEstore creates an empty store with room for ln records.
func newEstore(f byteid.IndexFactory, ln int) *estore {
	return &estore{
		idxFact:    f,
		bufferSize: ln,
		idx:        f(ln),
		records:    make([][]byte, ln),
	}
}

// factory keeps its stores in memory, indexed by name.
type factory struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	estores    []*estore
}

// FlatStore returns the store for name, creating an empty one the first time
// name is used. It never returns an error.
func (f *factory) FlatStore(bkt []byte) (store.FlatStore, error) {
	idx, found := f.idx.Get(bkt)
	if found {
		return f.estores[idx], nil
	}
	e := newEstore(f.idxFact, f.bufferSize)
	idx, app := f.idx.Insert(slices.Clone(bkt))
	if app {
		f.estores = append(f.estores, e)
	} else {
		f.estores[idx] = e
	}
	return e, nil
}

// Factory creates a store.FlatFactory that keeps every store in memory. Each
// store and the factory itself use an Index made by idxFact, and bufferSize is
// the length of the slice each one starts with, so it should be the number of
// records you expect. It is not safe for concurrent use.
func Factory(idxFact byteid.IndexFactory, bufferSize int) store.FlatFactory {
	return &factory{
		idxFact:    idxFact,
		idx:        idxFact(bufferSize),
		bufferSize: bufferSize,
		estores:    make([]*estore, bufferSize),
	}
}
