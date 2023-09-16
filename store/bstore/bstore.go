package bstore

import (
	"bytes"
	"os"

	"github.com/adamcolton/luce/store"
	"github.com/boltdb/bolt"
)

// boltstore is a store.NestedStore held in bolt buckets. bkt is the chain of
// bucket names: a top-level bucket followed by one nested bucket per level of
// Sub-Store. Its records are the bucket's keys and the buckets nested in it.
type boltstore struct {
	db  *bolt.DB
	bkt [][]byte
}

// Len counts the keys and Sub-Stores directly in the store. It takes O(n).
func (s *boltstore) Len() int {
	var ln int
	s.db.View(func(tx *bolt.Tx) error {
		if bkt := s.getBkt(tx); bkt != nil {
			c := bkt.Cursor()
			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				ln++
			}
		}
		return nil
	})
	return ln
}

// createBkt gets the store's bucket, creating it and any missing parent buckets.
func (s *boltstore) createBkt(tx *bolt.Tx) (*bolt.Bucket, error) {
	bkt, err := tx.CreateBucketIfNotExists(s.bkt[0])
	for i := 1; err == nil && i < len(s.bkt); i++ {
		bkt, err = bkt.CreateBucketIfNotExists(s.bkt[i])
	}
	return bkt, err
}

// getBkt gets the store's bucket, or nil if it or a parent bucket is missing.
func (s *boltstore) getBkt(tx *bolt.Tx) *bolt.Bucket {
	bkt := tx.Bucket(s.bkt[0])
	if bkt == nil {
		return nil
	}
	for _, bID := range s.bkt[1:] {
		bkt = bkt.Bucket(bID)
		if bkt == nil {
			return nil
		}
	}
	return bkt
}

// Put stores value at key, replacing any existing value. It returns an error if
// key holds a Sub-Store, or if the store's bucket can't be created because a
// bucket name it needs is now a key holding a value.
func (s *boltstore) Put(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt, err := s.createBkt(tx)
		if err != nil {
			return err
		}
		return bkt.Put(key, value)
	})
}

// Get returns the value stored at key, or the Sub-Store if key names one. The
// value is a copy. If the store's bucket is missing the key is not found.
func (s *boltstore) Get(key []byte) store.Record {
	var r store.Record
	s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getBkt(tx)
		if bkt == nil {
			return nil
		}
		if v := bkt.Get(key); v != nil {
			// v is only valid inside the transaction.
			r.Found = true
			r.Value = append(make([]byte, 0, len(v)), v...)
			return nil
		}
		if bkt.Bucket(key) != nil {
			r.Found = true
			r.Store = s.sub(key)
		}
		return nil
	})
	return r
}

// Next returns a copy of the smallest key greater than key, or of the lowest key
// if key is nil. It returns nil if there is no such key. Sub-Stores are keys.
func (s *boltstore) Next(key []byte) []byte {
	var nextKey []byte
	s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getBkt(tx)
		if bkt == nil {
			return nil
		}
		c := bkt.Cursor()
		k, _ := c.Seek(key)
		if k != nil && bytes.Equal(k, key) {
			k, _ = c.Next()
		}
		if k != nil {
			// k is only valid inside the transaction.
			nextKey = append(make([]byte, 0, len(k)), k...)
		}
		return nil
	})
	return nextKey
}

// Delete removes key, or the Sub-Store it names along with everything in it.
// Deleting a key that is not in the store is not an error.
func (s *boltstore) Delete(key []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt := s.getBkt(tx)
		if bkt == nil {
			return nil
		}
		bkt.Delete(key)
		bkt.DeleteBucket(key)
		return nil
	})
}

// NestedStore gets the Sub-Store at bkt, creating it if needed. It returns an
// error if bkt is empty or holds a value.
func (s *boltstore) NestedStore(bkt []byte) (store.NestedStore, error) {
	sub := s.sub(bkt)
	err := s.db.Update(func(tx *bolt.Tx) error {
		_, err := sub.createBkt(tx)
		return err
	})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// sub returns the Sub-Store at bkt without checking that it exists.
func (s *boltstore) sub(bkt []byte) *boltstore {
	ln := len(s.bkt)
	sub := &boltstore{
		db:  s.db,
		bkt: make([][]byte, ln+1),
	}
	copy(sub.bkt, s.bkt)
	sub.bkt[ln] = bkt
	return sub
}

// factory opens its bolt file when the first store is requested.
type factory struct {
	db          *bolt.DB
	permissions os.FileMode
	opts        *bolt.Options
	filename    string
}

// NestedStore returns the top-level store bkt. The store's bucket is created by
// its first Put. The first call opens the file. It returns an error if bkt is
// empty, because bolt needs a name for a bucket, or if the file can't be
// opened.
func (f *factory) NestedStore(bkt []byte) (store.NestedStore, error) {
	if len(bkt) == 0 {
		return nil, bolt.ErrBucketNameRequired
	}
	if f.db == nil {
		db, err := bolt.Open(f.filename, f.permissions, f.opts)
		if err != nil {
			return nil, err
		}
		f.db = db
	}
	return &boltstore{
		db:  f.db,
		bkt: [][]byte{bkt},
	}, nil
}

// Close closes the bolt file if it was opened. Stores made by the factory stop
// working.
func (f *factory) Close() error {
	if f.db == nil {
		return nil
	}
	return f.db.Close()
}

// FactoryCloser is a store.NestedFactory that holds open files, which Close
// releases.
type FactoryCloser interface {
	store.NestedFactory
	Close() error
}

// Factory creates a FactoryCloser for the bolt file at filename. Each store it
// returns is a top-level bucket in that file. The file is created and opened by
// the first call to NestedStore, using permissions and opts. It is opened once,
// and a bolt file can only be open once at a time.
func Factory(filename string, permissions os.FileMode, opts *bolt.Options) FactoryCloser {
	return &factory{
		filename:    filename,
		permissions: permissions,
		opts:        opts,
	}
}
