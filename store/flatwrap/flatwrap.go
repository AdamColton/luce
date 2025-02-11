// Package flatwrap builds a nested store on top of a flat store.
package flatwrap

import (
	"bytes"
	"sort"
	"sync/atomic"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store"
)

const (
	// ErrDeletedStore is returned by a store whose bucket has been deleted.
	ErrDeletedStore = lerr.Str("this bucket has been deleted")

	// ErrBktExists is returned when attempting to Put to a key that defines a
	// bucket.
	ErrBktExists = lerr.Str("Bucket already exists at that key")

	// ErrValExists is returned when attempting to create a Store with a key that
	// defines a value.
	ErrValExists = lerr.Str("Value already exists at that key")
)

// New creates a NestedFactory that keeps each store in the FlatStore of the same
// name. Every bucket in a store is a range of keys in the FlatStore: each key is
// stored with the 4 byte ID of its bucket in front of it. The IDs are stored in
// the FlatStore too, so a new factory over the same FlatFactory finds the
// buckets that an earlier one made. Nothing is safe for concurrent use, and keys
// must not be empty.
func New(flat store.FlatFactory) store.NestedFactory {
	return &factory{
		flat:  flat,
		roots: make(map[string]*virtualBucket),
	}
}

// factory hands out one store per name, so every caller that asks for the same
// name shares the state of its buckets.
type factory struct {
	flat  store.FlatFactory
	roots map[string]*virtualBucket
}

// NestedStore returns the store for name. Asking for the same name again
// returns the same store.
func (f *factory) NestedStore(name []byte) (store.NestedStore, error) {
	if vb, found := f.roots[string(name)]; found {
		return vb, nil
	}
	flat, err := f.flat.FlatStore(name)
	if err != nil {
		return nil, err
	}

	vb := &virtualBucket{
		root: &root{
			flat: flat,
		},
	}

	rec := flat.Get(maxIDKey)
	if rec.Found {
		vb.root.maxID = rye.Deserialize.Uint32(rec.Value)
	}

	vb.init()
	f.roots[string(name)] = vb
	return vb, nil
}

// maxIDKey holds the highest bucket ID handed out. It is shorter than any
// bucket key, which are at least 4 bytes.
var maxIDKey = []byte{0, 0}

// root is shared by every bucket in one store.
type root struct {
	flat  store.FlatStore
	maxID uint32
}

// virtualBucket is a store.NestedStore that is a range of keys in a FlatStore.
// The bucket's key at prefix alone holds the gob encoded names of its
// sub-buckets, and the key at prefix and a name holds the ID of that sub-bucket.
// Its other keys hold values. The root bucket has ID 0.
type virtualBucket struct {
	id       uint32
	prefix   [4]byte
	root     *root
	children map[string]*virtualBucket
	deleted  bool
}

// init loads the bucket's sub-buckets, or starts it empty if there are none.
func (vb *virtualBucket) init() {
	rye.Serialize.Uint32(vb.prefix[:], vb.id)
	rec := vb.root.flat.Get(vb.prefix[:])
	if rec.Found {
		vb.load(rec.Value)
	} else {
		vb.new()
	}
}

func (vb *virtualBucket) new() {
	vb.children = make(map[string]*virtualBucket)
}

// fullkey is the key in the FlatStore.
func (vb *virtualBucket) fullkey(key []byte) []byte {
	return append(vb.prefix[:], key...)
}

// load reads the sub-buckets from the encoded names. Each is created without
// its own sub-buckets, which are loaded when it is first used.
func (vb *virtualBucket) load(data []byte) {
	keys := []string{}
	gob.Dec(data, &keys)
	vb.children = make(map[string]*virtualBucket, len(keys))

	var buf []byte
	for _, key := range keys {
		ln := len(key) + 4
		if cap(buf) < ln {
			buf = make([]byte, ln*2)
			copy(buf, vb.prefix[:])
		}
		buf = buf[:ln]
		copy(buf[4:], key)
		rec := vb.root.flat.Get(buf)
		if rec.Found {
			vb.children[string(key)] = &virtualBucket{
				id:   rye.Deserialize.Uint32(rec.Value),
				root: vb.root,
			}
		}
	}
}

// saveChildren stores the names of the sub-buckets.
func (vb *virtualBucket) saveChildren() error {
	names := lmap.New(vb.children).Keys(nil)
	sort.Strings(names)
	return vb.root.flat.Put(vb.prefix[:], gob.Enc(names))
}

// Put stores a value. It returns ErrBktExists if key is a bucket and
// ErrDeletedStore if this bucket has been deleted.
func (vb *virtualBucket) Put(key, value []byte) error {
	if vb.deleted {
		return ErrDeletedStore
	}
	_, found := vb.children[string(key)]
	if found {
		return ErrBktExists
	}
	return vb.root.flat.Put(vb.fullkey(key), value)
}

// Get returns the value at key, or the bucket if key is one. A deleted bucket
// has nothing in it.
func (vb *virtualBucket) Get(key []byte) store.Record {
	if vb.deleted {
		return store.Record{}
	}
	child, found := vb.children[string(key)]
	if found {
		if child.children == nil {
			child.init()
		}
		return store.Record{
			Found: true,
			Store: child,
		}
	}
	return vb.root.flat.Get(vb.fullkey(key))
}

// owns is true if the FlatStore key is a value or bucket of this bucket.
func (vb *virtualBucket) owns(key []byte) bool {
	return len(key) > 4 && bytes.Equal(vb.prefix[:], key[:4])
}

// Next returns the smallest key greater than key, or the lowest key if key is
// nil. Buckets are keys. It returns nil if there is no such key.
func (vb *virtualBucket) Next(key []byte) (nextKey []byte) {
	next := vb.root.flat.Next(vb.fullkey(key))
	if vb.owns(next) {
		return next[4:]
	}
	return nil
}

// delete removes the bucket from the FlatStore along with every bucket in it,
// and marks them deleted.
func (vb *virtualBucket) delete() error {
	if vb.children == nil {
		vb.init()
	}
	vb.deleted = true
	for _, child := range vb.children {
		if err := child.delete(); err != nil {
			return err
		}
	}
	if err := vb.root.flat.Delete(vb.prefix[:]); err != nil {
		return err
	}
	for key := vb.Next(nil); key != nil; key = vb.Next(nil) {
		if err := vb.root.flat.Delete(vb.fullkey(key)); err != nil {
			return err
		}
	}
	return nil
}

// Delete removes a value, or a bucket and everything in it. Stores of the
// deleted bucket return ErrDeletedStore. Deleting a key that is not in the
// store is not an error.
func (vb *virtualBucket) Delete(key []byte) error {
	if vb.deleted {
		return ErrDeletedStore
	}
	c, found := vb.children[string(key)]
	if !found {
		return vb.root.flat.Delete(vb.fullkey(key))
	}
	if err := c.delete(); err != nil {
		return err
	}
	delete(vb.children, string(key))
	if err := vb.root.flat.Delete(vb.fullkey(key)); err != nil {
		return err
	}
	return vb.saveChildren()
}

// Len counts the values and buckets in the bucket by walking its keys, so it
// takes O(n).
func (vb *virtualBucket) Len() int {
	c := 0
	for k := vb.Next(nil); k != nil; k = vb.Next(k) {
		c++
	}
	return c
}

// NestedStore gets the bucket at key, creating it if needed. It returns
// ErrValExists if key is a value and ErrDeletedStore if this bucket has been
// deleted. A bucket that can't be saved is not created.
func (w *virtualBucket) NestedStore(key []byte) (store.NestedStore, error) {
	if w.deleted {
		return nil, ErrDeletedStore
	}
	check := w.Get(key)
	if check.Found {
		if check.Store == nil {
			return nil, ErrValExists
		}
		return check.Store, nil
	}

	c := &virtualBucket{
		id:   atomic.AddUint32(&(w.root.maxID), 1),
		root: w.root,
	}

	// The highest ID is saved first, so a failure can't leave an ID that a
	// later bucket would use again.
	var mp [4]byte
	rye.Serialize.Uint32(mp[:], w.root.maxID)
	if err := w.root.flat.Put(maxIDKey, mp[:]); err != nil {
		return nil, err
	}

	c.init()
	if err := w.root.flat.Put(w.fullkey(key), c.prefix[:]); err != nil {
		return nil, err
	}
	w.children[string(key)] = c
	if err := w.saveChildren(); err != nil {
		delete(w.children, string(key))
		w.root.flat.Delete(w.fullkey(key))
		return nil, err
	}
	return c, nil
}
