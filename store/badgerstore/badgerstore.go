package badgerstore

import (
	"path"

	"github.com/adamcolton/luce/store"
	badger "github.com/dgraph-io/badger/v4"
)

type badgerStore struct {
	db *badger.DB
}

func (s *badgerStore) Put(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *badgerStore) Get(key []byte) store.Record {
	var r store.Record
	s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			r.Found = true
			r.Value = val
			return nil
		})
	})
	return r
}

func (s *badgerStore) Next(key []byte) (nextKey []byte) {
	s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		it.Seek(key)
		if it.Valid() {
			it.Next()
			if it.Valid() {
				nextKey = it.Item().Key()
			}
		}
		it.Close()
		return nil
	})
	return
}

func (s *badgerStore) Delete(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (s *badgerStore) Close() error {
	return s.db.Close()
}

func (s *badgerStore) Sync() error {
	return s.db.Sync()
}

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

type factory struct {
	root string
	dbs  map[string]*badgerStore
}

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

func Factory(root string) store.FlatFactory {
	return &factory{
		root: root,
		dbs:  make(map[string]*badgerStore),
	}
}
