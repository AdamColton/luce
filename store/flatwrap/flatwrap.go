package flatwrap

import (
	"bytes"
	"sync/atomic"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store"
)

func New(flat store.FlatFactory) store.NestedFactory {
	return factory{
		flat: flat,
	}
}

type factory struct {
	flat store.FlatFactory
}

func (f factory) NestedStore(name []byte) (store.NestedStore, error) {
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
	return vb, nil
}

var maxIDKey = []byte{0, 0}

type root struct {
	flat  store.FlatStore
	maxID uint32
}

type virtualBucket struct {
	id       uint32
	prefix   [4]byte
	root     *root
	children map[string]*virtualBucket
	deleted  bool
}

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

func (vb *virtualBucket) fullkey(key []byte) []byte {
	return append(vb.prefix[:], key...)
}

// at the key for a sub-virtual bucket we store the id

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

const ErrDeletedStore = lerr.Str("this bucket has been deleted")

func (vb *virtualBucket) Put(key, value []byte) error {
	if vb.deleted {
		return ErrDeletedStore
	}
	_, found := vb.children[string(key)]
	if found {
		return lerr.Str("bucket already exists at that key")
	}
	return vb.root.flat.Put(vb.fullkey(key), value)
}

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

func (vb *virtualBucket) owns(key []byte) bool {
	return len(key) > 4 && bytes.Equal(vb.prefix[:], key[:4])
}

func (vb *virtualBucket) Next(key []byte) (nextKey []byte) {
	next := vb.root.flat.Next(vb.fullkey(key))
	if vb.owns(next) {
		return next[4:]
	}
	return nil
}

func (vb *virtualBucket) delete() {
	if vb.children == nil {
		vb.init()
	}
	for _, child := range vb.children {
		child.delete()
	}
	key := vb.prefix[:]
	for {
		vb.root.flat.Delete(key)
		key = vb.Next(key)
		if !vb.owns(key) {
			break
		}
	}
}

func (vb *virtualBucket) Delete(key []byte) error {
	if vb.deleted {
		return ErrDeletedStore
	}
	c, found := vb.children[string(key)]
	if found {
		c.delete()
		delete(vb.children, string(key))
	}
	return vb.root.flat.Delete(vb.fullkey(key))
}

func (vb *virtualBucket) Len() int {
	c := 0
	for k := vb.Next(nil); k != nil; k = vb.Next(k) {
		c++
	}
	return c
}

func (w *virtualBucket) NestedStore(key []byte) (store.NestedStore, error) {
	check := w.Get(key)
	if check.Found {
		if check.Store == nil {
			return nil, lerr.Str("Value already exists at that key")
		}
		return check.Store, nil
	}

	c := &virtualBucket{
		id:   atomic.AddUint32(&(w.root.maxID), 1),
		root: w.root,
	}
	c.init()
	w.children[string(key)] = c
	w.root.flat.Put(w.fullkey(key), c.prefix[:])

	var childIDs []string = lmap.New(w.children).Keys(nil)
	data := gob.Enc(childIDs)
	w.root.flat.Put(w.prefix[:], data)

	var mp [4]byte
	rye.Serialize.Uint32(mp[:], w.root.maxID)
	w.root.flat.Put(maxIDKey, mp[:])

	return c, nil
}
