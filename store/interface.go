package store

// FlatStore is a store of key/value pairs ordered by key. Some implementations
// keep the slices they are given, so don't modify a key or value after a Put.
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

// NestedStore is a FlatStore that can also hold named Sub-Stores. A key holds
// either a value or a Sub-Store, never both.
type NestedStore interface {
	FlatStore
	// NestedFactory allows Sub-Stores to be created.
	NestedFactory
}

// NestedFactory creates and gets NestedStores by name.
type NestedFactory interface {
	// NestedStore acts as an Upsert operation, it will get the Store if it exists
	// and create it if it does not.
	NestedStore(name []byte) (NestedStore, error)
}

// FlatFactory creates and gets FlatStores by name.
type FlatFactory interface {
	// FlatStore acts as an Upsert operation, it will get the Store if it exists
	// and create it if it does not.
	FlatStore(name []byte) (FlatStore, error)
}

// Record is returned from a call to FlatStore.Get. Found indicates if the key
// was found. If it was, Value holds the value, or Store holds the Sub-Store if
// the key names one.
type Record struct {
	Found bool
	Value []byte
	Store NestedStore
}

// Flattener turns a NestedFactory into a FlatFactory. The FlatStore it returns
// is the NestedStore of the same name.
type Flattener struct {
	NestedFactory
}

// FlatStore fulfills FlatFactory. It returns the NestedStore that the
// NestedFactory holds under name, along with any error the factory returns.
func (f Flattener) FlatStore(name []byte) (FlatStore, error) {
	ns, err := f.NestedFactory.NestedStore(name)
	return ns, err
}
