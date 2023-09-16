package store

// Store of key/value pairs
type Store interface {
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

	// Store fulfills Factory allowing Sub-Stores to be created.
	Factory
}

// Factory for creating a Store
type Factory interface {
	// Store acts as an Upsert operation, it will get the Store if it exists
	// and create it if it does not.
	Store(name []byte) (Store, error)
}

// Record is returned from a call to Store.Get. Found indicates if the key was
// found. If it was either the Value or the Store will not be nil.
type Record struct {
	Found bool
	Value []byte
	Store Store
}
