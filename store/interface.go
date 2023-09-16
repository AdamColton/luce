package store

// Store of key/value pairs
type FlatStore interface {
	// Put a key,value pair. This will over write an existing value. An error
	// is returned if there is a collision with a Sub-Store.
	Put(key, value []byte) error
	// Get returns the value stored at the key. The Record will contain either
	// a Value or a Store if the key exists.
	Get(key []byte) Record
	// Next returns the next key greater than the one provided. If nil is passed
	// in, the lowest key is returned. If the highest key in the store is passed
	// in, nil is returned.
	Next(key []byte) (nextKey []byte)
	// Delete will delete either a key or a Sub-Store
	Delete(key []byte) error
	// Fulfills slice.Lener, returns how many records are stored.
	Len() int
}

type NestedStore interface {
	FlatStore
	// Store fulfills Factory allowing Sub-Stores to be created.
	NestedFactory
}

// Factory for creating a Store
type NestedFactory interface {
	// Store acts as an Upsert operation, it will get the Store if it exists
	// and create it if it does not.
	NestedStore(name []byte) (NestedStore, error)
}

type FlatFactory interface {
	// Store acts as an Upsert operation, it will get the Store if it exists
	// and create it if it does not.
	FlatStore(name []byte) (FlatStore, error)
}

// Record is returned from a call to Store.Get. Found indicates if the key was
// found. If it was either the Value or the Store will not be nil.
type Record struct {
	Found bool
	Value []byte
	Store NestedStore
}

type Flattener struct {
	NestedFactory
}

func (f Flattener) FlatStore(name []byte) (FlatStore, error) {
	ns, err := f.NestedFactory.NestedStore(name)
	return ns, err
}
