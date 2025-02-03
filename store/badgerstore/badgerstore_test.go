package badgerstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/badgerstore"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

// closer is implemented by the stores that badgerstore returns.
type closer interface {
	Close() error
	Sync() error
}

// factory creates a Factory in a temp directory and closes the stores it made.
func factory(t *testing.T, names ...string) store.FlatFactory {
	f := badgerstore.Factory(t.TempDir())
	t.Cleanup(func() {
		for _, n := range names {
			s, _ := f.FlatStore([]byte(n))
			s.(closer).Close()
		}
	})
	return f
}

func TestBasic(t *testing.T) {
	testsuite.TestFlat(t, factory(t, "TestBasic"))
}

func TestIteration(t *testing.T) {
	testsuite.TestIteration(t, factory(t, "TestIteration"))
}

func TestNext(t *testing.T) {
	s, err := factory(t, "next").FlatStore([]byte("next"))
	assert.NoError(t, err)
	assert.Nil(t, s.Next(nil))

	for _, k := range []string{"b", "d", "f"} {
		assert.NoError(t, s.Put([]byte(k), []byte(k+k)))
	}
	assert.Equal(t, []byte("b"), s.Next(nil))
	// The key doesn't have to be in the store.
	assert.Equal(t, []byte("b"), s.Next([]byte("a")))
	assert.Equal(t, []byte("d"), s.Next([]byte("b")))
	assert.Equal(t, []byte("f"), s.Next([]byte("e")))
	assert.Nil(t, s.Next([]byte("f")))
	assert.Nil(t, s.Next([]byte("g")))
}

func TestGetDelete(t *testing.T) {
	s, err := factory(t, "get").FlatStore([]byte("get"))
	assert.NoError(t, err)

	assert.Equal(t, store.Record{}, s.Get([]byte("apple")))
	assert.NoError(t, s.Put([]byte("apple"), []byte("red")))
	assert.NoError(t, s.Put([]byte("apple"), []byte("green")))
	assert.Equal(t, store.Record{Found: true, Value: []byte("green")}, s.Get([]byte("apple")))
	assert.Equal(t, 1, s.Len())

	assert.NoError(t, s.Delete([]byte("apple")))
	assert.NoError(t, s.Delete([]byte("apple")))
	assert.False(t, s.Get([]byte("apple")).Found)
	assert.Equal(t, 0, s.Len())
}

func TestFactory(t *testing.T) {
	f := factory(t, "same")
	a, err := f.FlatStore([]byte("same"))
	assert.NoError(t, err)
	b, err := f.FlatStore([]byte("same"))
	assert.NoError(t, err)
	assert.Same(t, a, b)

	// Each name is a separate store.
	c, err := f.FlatStore([]byte("other"))
	assert.NoError(t, err)
	defer c.(closer).Close()
	assert.NoError(t, a.Put([]byte("k"), []byte("v")))
	assert.False(t, c.Get([]byte("k")).Found)
}

func TestFactoryError(t *testing.T) {
	// The database directory can't be created below a file.
	file := filepath.Join(t.TempDir(), "file")
	assert.NoError(t, os.WriteFile(file, nil, 0600))

	s, err := badgerstore.Factory(file).FlatStore([]byte("name"))
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestReopen(t *testing.T) {
	root := t.TempDir()
	name := []byte("reopen")
	s, err := badgerstore.Factory(root).FlatStore(name)
	assert.NoError(t, err)
	assert.NoError(t, s.Put([]byte("k"), []byte("v")))
	assert.NoError(t, s.(closer).Sync())
	assert.NoError(t, s.(closer).Close())

	s, err = badgerstore.Factory(root).FlatStore(name)
	assert.NoError(t, err)
	defer s.(closer).Close()
	assert.Equal(t, []byte("v"), s.Get([]byte("k")).Value)
}
