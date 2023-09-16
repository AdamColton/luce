package bstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/bstore"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

// factory creates a Factory on a file in a temp directory and closes it when
// the test ends.
func factory(t *testing.T) bstore.FactoryCloser {
	f := bstore.Factory(filepath.Join(t.TempDir(), "test.db"), 0600, nil)
	t.Cleanup(func() { assert.NoError(t, f.Close()) })
	return f
}

func TestBasic(t *testing.T) {
	testsuite.TestAll(t, factory(t))
}

func TestEmptyName(t *testing.T) {
	f := factory(t)

	// Bolt needs a name for every bucket.
	s, err := f.NestedStore(nil)
	assert.Equal(t, bolt.ErrBucketNameRequired, err)
	assert.Nil(t, s)

	s, err = f.NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub, err := s.NestedStore([]byte{})
	assert.Error(t, err)
	assert.Nil(t, sub)
}

func TestOpenError(t *testing.T) {
	// The file can't be created in a directory that doesn't exist.
	f := bstore.Factory(filepath.Join(t.TempDir(), "missing", "test.db"), 0600, nil)
	s, err := f.NestedStore([]byte("root"))
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.NoError(t, f.Close())
}

func TestReopen(t *testing.T) {
	name := filepath.Join(t.TempDir(), "test.db")
	f := bstore.Factory(name, 0600, nil)
	s, err := f.NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub, err := s.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	assert.NoError(t, sub.Put([]byte("k"), []byte("v")))
	assert.NoError(t, f.Close())

	f = bstore.Factory(name, 0600, nil)
	defer f.Close()
	s, err = f.NestedStore([]byte("root"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("v"), s.Get([]byte("sub")).Store.Get([]byte("k")).Value)
}

func TestValues(t *testing.T) {
	s, err := factory(t).NestedStore([]byte("values"))
	assert.NoError(t, err)

	// An empty value is still found.
	assert.NoError(t, s.Put([]byte("empty"), nil))
	assert.Equal(t, store.Record{Found: true, Value: []byte{}}, s.Get([]byte("empty")))
	assert.Equal(t, store.Record{}, s.Get([]byte("missing")))

	// What Get and Next return can be modified without changing the store.
	assert.NoError(t, s.Put([]byte("key"), []byte("value")))
	s.Get([]byte("key")).Value[0] = 'V'
	assert.Equal(t, []byte("value"), s.Get([]byte("key")).Value)
	s.Next([]byte("empty"))[0] = 'K'
	assert.Equal(t, []byte("key"), s.Next([]byte("empty")))

	assert.Equal(t, []byte("empty"), s.Next(nil))
	assert.Nil(t, s.Next([]byte("key")))
	assert.NoError(t, s.Delete([]byte("key")))
	assert.NoError(t, s.Delete([]byte("key")))
	assert.False(t, s.Get([]byte("key")).Found)
}

func TestLen(t *testing.T) {
	s, err := factory(t).NestedStore([]byte("len"))
	assert.NoError(t, err)
	assert.Equal(t, 0, s.Len())

	assert.NoError(t, s.Put([]byte("a"), []byte("1")))
	sub, err := s.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	assert.NoError(t, sub.Put([]byte("x"), []byte("1")))
	assert.NoError(t, sub.Put([]byte("y"), []byte("2")))

	// A Sub-Store is one record, and its own records aren't counted.
	assert.Equal(t, 2, s.Len())
	assert.Equal(t, 2, sub.Len())
}

func TestStaleStore(t *testing.T) {
	root, err := factory(t).NestedStore([]byte("stale"))
	assert.NoError(t, err)
	sub, err := root.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	assert.NoError(t, sub.Put([]byte("x"), []byte("1")))

	// The Sub-Store is deleted, so the store is empty.
	assert.NoError(t, root.Delete([]byte("sub")))
	assert.False(t, sub.Get([]byte("x")).Found)
	assert.Nil(t, sub.Next(nil))
	assert.Equal(t, 0, sub.Len())
	assert.NoError(t, sub.Delete([]byte("x")))

	// Now its name is a value, so it can't be created again.
	assert.NoError(t, root.Put([]byte("sub"), []byte("value")))
	assert.Error(t, sub.Put([]byte("y"), []byte("2")))
	assert.Equal(t, []byte("value"), root.Get([]byte("sub")).Value)
}

func TestDirectory(t *testing.T) {
	dir := t.TempDir()
	d := bstore.Directory(dir, 0600, nil)
	defer d.Close()

	fruit, err := d.NestedStore([]byte("fruit"))
	assert.NoError(t, err)
	animals, err := d.NestedStore([]byte("animals"))
	assert.NoError(t, err)

	assert.NoError(t, fruit.Put([]byte("apple"), []byte("red")))
	assert.NoError(t, animals.Put([]byte("cat"), []byte("meow")))
	assert.Equal(t, []byte("red"), fruit.Get([]byte("apple")).Value)
	assert.False(t, animals.Get([]byte("apple")).Found)
	assert.Equal(t, 1, fruit.Len())

	// Each store is a file in the directory.
	for _, name := range []string{"fruit", "animals"} {
		_, err := os.Stat(filepath.Join(dir, name))
		assert.NoError(t, err)
	}

	// The same name is the same store.
	again, err := d.NestedStore([]byte("fruit"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("red"), again.Get([]byte("apple")).Value)

	// Sub-Stores work as they do elsewhere.
	sub, err := fruit.NestedStore([]byte("citrus"))
	assert.NoError(t, err)
	assert.NoError(t, sub.Put([]byte("lemon"), []byte("sour")))

	// Closing releases the files, and the data is still there.
	assert.NoError(t, d.Close())
	d = bstore.Directory(dir, 0600, nil)
	fruit, err = d.NestedStore([]byte("fruit"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("sour"), fruit.Get([]byte("citrus")).Store.Get([]byte("lemon")).Value)
	assert.Equal(t, 2, fruit.Len())
}

func TestDirectoryOpenError(t *testing.T) {
	d := bstore.Directory(filepath.Join(t.TempDir(), "missing"), 0600, nil)
	s, err := d.NestedStore([]byte("fruit"))
	assert.Error(t, err)
	assert.Nil(t, s)
	assert.NoError(t, d.Close())
}
