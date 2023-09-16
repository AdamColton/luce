package bstore

import (
	"bytes"
	"os"

	"github.com/adamcolton/luce/store"
	"github.com/boltdb/bolt"
)

type boltstore struct {
	db  *bolt.DB
	bkt [][]byte
}

func (s *boltstore) Len() int {
	var ln int
	s.db.View(func(tx *bolt.Tx) error {
		if bkt := s.getBkt(tx); bkt != nil {
			ln = bkt.Stats().KeyN
		}
		return nil
	})
	return ln
}

func (s *boltstore) createBkt(tx *bolt.Tx) (*bolt.Bucket, error) {
	bkt, err := tx.CreateBucketIfNotExists(s.bkt[0])
	if bkt == nil || err != nil {
		return bkt, err
	}
	for _, bID := range s.bkt[1:] {
		bkt, err = bkt.CreateBucketIfNotExists(bID)
		if err != nil {
			return nil, err
		}
	}
	return bkt, nil
}
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

func (s *boltstore) Put(key, value []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bkt, err := s.createBkt(tx)
		if err != nil {
			return nil
		}
		return bkt.Put(key, value)
	})
}

func (s *boltstore) Get(key []byte) store.Record {
	var r store.Record
	s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getBkt(tx)
		if bkt == nil {
			return nil
		}
		r.Value = bkt.Get(key)
		if r.Value != nil {
			r.Found = true
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

func (s *boltstore) Next(key []byte) []byte {
	var nextKey []byte
	s.db.View(func(tx *bolt.Tx) error {
		bkt := s.getBkt(tx)
		if bkt == nil {
			return nil
		}
		c := bkt.Cursor()
		nextKey, _ = c.Seek(key)
		if nextKey == nil {
			return nil
		}
		if bytes.Equal(nextKey, key) {
			nextKey, _ = c.Next()
		}
		return nil
	})
	return nextKey
}

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

type factory struct {
	db          *bolt.DB
	permissions os.FileMode
	opts        *bolt.Options
	filename    string
}

// Store returns the bolt store backed instance of store.
func (f *factory) NestedStore(bkt []byte) (store.NestedStore, error) {
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

func (f *factory) Close() error {
	if f.db == nil {
		return nil
	}
	return f.db.Close()
}

type FactoryCloser interface {
	store.NestedFactory
	Close() error
}

// Factory defines an instance of a bolt store.
func Factory(filename string, permissions os.FileMode, opts *bolt.Options) FactoryCloser {
	return &factory{
		filename:    filename,
		permissions: permissions,
		opts:        opts,
	}
}
