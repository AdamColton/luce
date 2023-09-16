package bstore

import (
	"os"
	"path"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
	"github.com/boltdb/bolt"
)

// rootBucket is the bucket that holds the data of a store made by Directory.
var rootBucket = []byte("root")

// directory holds the files it has opened, by store name.
type directory struct {
	permissions os.FileMode
	opts        *bolt.Options
	dir         string
	files       map[string]*factory
}

// Directory creates a FactoryCloser that keeps each store in its own bolt file,
// named for the store, in dir, which must exist. Each file holds the store's
// data in one root bucket. Files stay open until Close.
func Directory(dir string, permissions os.FileMode, opts *bolt.Options) FactoryCloser {
	return &directory{
		dir:         dir,
		permissions: permissions,
		opts:        opts,
		files:       make(map[string]*factory),
	}
}

// NestedStore returns the store for name, creating its file if needed. Asking
// for the same name again returns a store on the same open file. The name is
// used as a file name, so it should not contain a path separator.
func (d *directory) NestedStore(name []byte) (store.NestedStore, error) {
	f, found := d.files[string(name)]
	if !found {
		f = &factory{
			filename:    path.Join(d.dir, string(name)),
			permissions: d.permissions,
			opts:        d.opts,
		}
	}
	s, err := f.NestedStore(rootBucket)
	if err != nil {
		return nil, err
	}
	d.files[string(name)] = f
	return s, nil
}

// Close closes every file the directory opened. Its stores stop working, and
// asking for a name again opens its file again.
func (d *directory) Close() error {
	var errs lerr.Many
	for name, f := range d.files {
		errs = errs.Add(f.Close())
		delete(d.files, name)
	}
	return errs.Cast()
}
