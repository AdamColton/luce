package bstore

import (
	"os"
	"path"

	"github.com/adamcolton/luce/store"
	"github.com/boltdb/bolt"
)

type directory struct {
	permissions os.FileMode
	opts        *bolt.Options
	dir         string
}

// Directory creates a Factory that will create a new bolt file for each store
// that is opened.
func Directory(dir string, permissions os.FileMode, opts *bolt.Options) store.Factory {
	return &directory{
		dir:         dir,
		permissions: permissions,
		opts:        opts,
	}
}

// Store creates a bolt file.
func (d *directory) Store(bkt []byte) (store.Store, error) {
	return Factory(path.Join(d.dir, string(bkt)), d.permissions, d.opts).Store(nil)
}
