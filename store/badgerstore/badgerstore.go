// Package badgerstore implements the store interfaces with BadgerDB.
package badgerstore

import (
	"bytes"
	"path"

	"github.com/adamcolton/luce/store"
	badger "github.com/dgraph-io/badger/v4"
)

// badgerStore is a store.FlatStore held in one BadgerDB database.
type badgerStore struct {
	db *badger.DB
}

// Put stores value at key, replacing any existing value.
func (s *badgerStore) Put(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Get returns a copy of the value stored at key. A read error is reported as a
// key that was not found.
func (s *badgerStore) Get(key []byte) store.Record {
	var r store.Record
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			// val is only valid inside this func.
			r.Found = true
			r.Value = append([]byte(nil), val...)
			return nil
		})
	})
	return r
}

// Next returns a copy of the smallest key greater than key, or of the lowest key
// if key is nil. It returns nil if there is no such key or on a read error.
func (s *badgerStore) Next(key []byte) (nextKey []byte) {
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		it.Seek(key)
		if it.Valid() && bytes.Equal(it.Item().Key(), key) {
			it.Next()
		}
		if it.Valid() {
			nextKey = it.Item().KeyCopy(nil)
		}
		it.Close()
		return nil
	})
	return
}

// Delete removes key. Deleting a key that is not in the store is not an error.
func (s *badgerStore) Delete(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Close closes the database. The factory keeps returning the closed store for
// its name, so the name can't be reopened through the same factory.
func (s *badgerStore) Close() error {
	return s.db.Close()
}

// Sync flushes writes to disk.
func (s *badgerStore) Sync() error {
	return s.db.Sync()
}

// Len counts the keys by iterating over all of them, so it takes O(n). A read
// error ends the count early.
func (s *badgerStore) Len() (ln int) {
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		for it.Rewind(); it.Valid(); it.Next() {
			ln++
		}
		it.Close()
		return nil
	})
	return
}

// factory holds the databases it has opened, by name.
type factory struct {
	root string
	dbs  map[string]*badgerStore
}

// FlatStore returns the store for name. The first call for a name opens the
// database in the directory root/name, creating it if needed, and later calls
// return that same store.
func (f *factory) FlatStore(name []byte) (store.FlatStore, error) {
	ns := string(name)
	b, found := f.dbs[ns]
	if found {
		return b, nil
	}

	p := path.Join(f.root, ns)
	opts := badger.DefaultOptions(p)
	opts.Logger = nil
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	b = &badgerStore{
		db: db,
	}
	f.dbs[ns] = b
	return b, nil
}

// Factory returns a store.FlatFactory that keeps each store in its own BadgerDB
// database, in a directory named for the store under root. The stores also have
// Close and Sync methods. The factory is not safe for concurrent use.
func Factory(root string) store.FlatFactory {
	return &factory{
		root: root,
		dbs:  make(map[string]*badgerStore),
	}
}
