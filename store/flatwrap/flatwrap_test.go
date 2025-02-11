package flatwrap_test

import (
	"errors"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/store/flatwrap"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	fac := flatwrap.New(ephemeral.Factory(bytebtree.New, 1))
	testsuite.TestAll(t, fac)
}

func TestReload(t *testing.T) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	fac := flatwrap.New(flat)
	fruit, err := fac.NestedStore([]byte("fruit"))
	assert.NoError(t, err)

	fmap := map[string]string{
		"A": "apple",
		"B": "banana",
		"C": "cantaloup",
		"D": "date",
		"E": "elderberry",
	}
	for k, v := range fmap {
		fruit.Put([]byte(k), []byte(v))
	}

	animal, err := fac.NestedStore([]byte("animal"))
	assert.NoError(t, err)

	amap := map[string]string{
		"A": "armadillo",
		"B": "badger",
		"C": "cat",
		"D": "dog",
		"E": "elephant",
	}
	for k, v := range amap {
		animal.Put([]byte(k), []byte(v))
	}

	fac2 := flatwrap.New(flat)
	fruit2, err := fac2.NestedStore([]byte("fruit"))
	assert.NoError(t, err)
	it := store.NewIter(fruit2)
	got := make(map[string]string)
	for !it.Done() {
		k, r, _ := it.CurVal()
		got[string(k)] = string(r.Value)
		it.Next()
	}
	assert.Equal(t, fmap, got)

	animal2, err := fac2.NestedStore([]byte("animal"))
	assert.NoError(t, err)
	it = store.NewIter(animal2)
	got = make(map[string]string)
	for !it.Done() {
		k, r, _ := it.CurVal()
		got[string(k)] = string(r.Value)
		it.Next()
	}
	assert.Equal(t, amap, got)
}

func TestNext(t *testing.T) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	root, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)

	root.Put([]byte("A"), []byte("A"))
	root.NestedStore([]byte("B"))
	root.Put([]byte("C"), []byte("C"))
	root.NestedStore([]byte("D"))
	root.Put([]byte("E"), []byte("E"))

	root2, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	got := make([]string, 0, 5)
	liter.Wrap(store.NewIter(root2)).For(func(key []byte) {
		got = append(got, string(key))
	})

	expected := []string{"A", "B", "C", "D", "E"}
	assert.Equal(t, expected, got)
}

// failing is a FlatStore that can be told to fail its nth Put or Delete, counting
// from when it was armed.
type failing struct {
	store.FlatStore
	putAt, deleteAt int
}

var errFailed = errors.New("failed")

func (f *failing) Put(key, value []byte) error {
	if f.putAt--; f.putAt == 0 {
		return errFailed
	}
	return f.FlatStore.Put(key, value)
}

func (f *failing) Delete(key []byte) error {
	if f.deleteAt--; f.deleteAt == 0 {
		return errFailed
	}
	return f.FlatStore.Delete(key)
}

// failingFactory makes the same failing store for every name, and can fail.
type failingFactory struct {
	store.FlatFactory
	fail *failing
	err  error
}

func (f failingFactory) FlatStore(name []byte) (store.FlatStore, error) {
	if f.err != nil {
		return nil, f.err
	}
	s, err := f.FlatFactory.FlatStore(name)
	f.fail.FlatStore = s
	return f.fail, err
}

// setup creates a root store with a value and a sub-store that has a value and
// a sub-store of its own. It returns the FlatStore behind them.
func setup(t *testing.T, fs *failing) (root, sub, inner store.NestedStore, raw store.FlatStore) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	var fac store.NestedFactory = flatwrap.New(flat)
	if fs != nil {
		fac = flatwrap.New(failingFactory{FlatFactory: flat, fail: fs})
	}
	var err error
	root, err = fac.NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub, err = root.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	inner, err = sub.NestedStore([]byte("inner"))
	assert.NoError(t, err)
	assert.NoError(t, root.Put([]byte("keep"), []byte("kept")))
	assert.NoError(t, sub.Put([]byte("a"), []byte("1")))
	assert.NoError(t, sub.Put([]byte("b"), []byte("2")))
	assert.NoError(t, inner.Put([]byte("c"), []byte("3")))
	raw, err = flat.FlatStore([]byte("root"))
	assert.NoError(t, err)
	return
}

func TestDelete(t *testing.T) {
	root, sub, inner, raw := setup(t, nil)
	// The FlatStore holds the highest ID, root's names list, its bucket and its
	// value (4), then sub's names list, its bucket, two values (4), and inner's
	// value (1).
	assert.Equal(t, 9, raw.Len())

	assert.NoError(t, root.Delete([]byte("sub")))
	// What is left is the highest ID, root's names list and its value.
	assert.Equal(t, 3, raw.Len())
	assert.Equal(t, 1, root.Len())
	assert.False(t, root.Get([]byte("sub")).Found)
	assert.NoError(t, root.Delete([]byte("sub")))
	assert.NoError(t, root.Delete([]byte("missing")))

	// Every store in the deleted bucket is deleted.
	for _, s := range []store.NestedStore{sub, inner} {
		assert.Equal(t, flatwrap.ErrDeletedStore, s.Put([]byte("x"), []byte("y")))
		assert.Equal(t, flatwrap.ErrDeletedStore, s.Delete([]byte("a")))
		_, err := s.NestedStore([]byte("x"))
		assert.Equal(t, flatwrap.ErrDeletedStore, err)
		assert.Equal(t, store.Record{}, s.Get([]byte("a")))
		assert.Nil(t, s.Next(nil))
		assert.Equal(t, 0, s.Len())
	}
	assert.Equal(t, 3, raw.Len())
}

func TestReloadNested(t *testing.T) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	root, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub, err := root.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	inner, err := sub.NestedStore([]byte("inner"))
	assert.NoError(t, err)
	assert.NoError(t, sub.Put([]byte("b"), []byte("2")))
	assert.NoError(t, inner.Put([]byte("c"), []byte("3")))

	// A new factory loads each bucket when it is first used.
	root2, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub2 := root2.Get([]byte("sub")).Store
	assert.Equal(t, []byte("2"), sub2.Get([]byte("b")).Value)
	assert.Equal(t, 2, sub2.Len())
	inner2 := sub2.Get([]byte("inner")).Store
	assert.Equal(t, []byte("3"), inner2.Get([]byte("c")).Value)

	// The IDs carry on from where they were.
	other, err := root2.NestedStore([]byte("other"))
	assert.NoError(t, err)
	assert.NoError(t, other.Put([]byte("c"), []byte("not 3")))
	assert.Equal(t, []byte("3"), inner2.Get([]byte("c")).Value)
}

func TestDeleteReloaded(t *testing.T) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	root, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	sub, err := root.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	inner, err := sub.NestedStore([]byte("inner"))
	assert.NoError(t, err)
	assert.NoError(t, inner.Put([]byte("c"), []byte("3")))
	assert.NoError(t, root.Put([]byte("other"), []byte("o")))

	// The buckets are not loaded until they are used, and the names of the
	// deleted ones are forgotten.
	root2, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	assert.NoError(t, root2.Delete([]byte("sub")))
	assert.NoError(t, root2.Put([]byte("sub"), []byte("now a value")))

	root3, err := flatwrap.New(flat).NestedStore([]byte("root"))
	assert.NoError(t, err)
	assert.Equal(t, store.Record{Found: true, Value: []byte("now a value")}, root3.Get([]byte("sub")))
	assert.Equal(t, 2, root3.Len())
	raw, _ := flat.FlatStore([]byte("root"))
	// The highest ID, root's names list, "other" and "sub".
	assert.Equal(t, 4, raw.Len())
}

func TestExists(t *testing.T) {
	root, sub, _, _ := setup(t, nil)

	// A key is a value or a bucket, not both.
	assert.Equal(t, flatwrap.ErrBktExists, root.Put([]byte("sub"), []byte("value")))
	_, err := root.NestedStore([]byte("keep"))
	assert.Equal(t, flatwrap.ErrValExists, err)
	assert.Equal(t, []byte("kept"), root.Get([]byte("keep")).Value)

	// A bucket that exists is returned.
	got, err := root.NestedStore([]byte("sub"))
	assert.NoError(t, err)
	assert.Same(t, sub, got)
}

func TestSameName(t *testing.T) {
	flat := ephemeral.Factory(bytebtree.New, 1)
	fac := flatwrap.New(flat)
	a, err := fac.NestedStore([]byte("x"))
	assert.NoError(t, err)
	b, err := fac.NestedStore([]byte("x"))
	assert.NoError(t, err)
	assert.Same(t, a, b)

	// Buckets made through either are seen by both, and don't share an ID.
	one, _ := a.NestedStore([]byte("one"))
	two, _ := b.NestedStore([]byte("two"))
	assert.NoError(t, one.Put([]byte("k"), []byte("from one")))
	assert.False(t, two.Get([]byte("k")).Found)
	assert.True(t, a.Get([]byte("two")).Found)
	assert.True(t, b.Get([]byte("one")).Found)

	// Another name is another store.
	c, err := fac.NestedStore([]byte("y"))
	assert.NoError(t, err)
	assert.NotSame(t, a, c)
	assert.Equal(t, 0, c.Len())
}

func TestFactoryError(t *testing.T) {
	fac := flatwrap.New(failingFactory{err: errFailed})
	s, err := fac.NestedStore([]byte("x"))
	assert.Equal(t, errFailed, err)
	assert.Nil(t, s)
}

func TestNestedStoreErrors(t *testing.T) {
	// Fail each Put in turn until NestedStore works.
	for n := 1; ; n++ {
		fs := &failing{}
		root, sub, _, raw := setup(t, fs)
		before := raw.Len()
		fs.putAt = n
		got, err := root.NestedStore([]byte("new"))
		if err == nil {
			assert.NotNil(t, got)
			assert.Greater(t, n, 1)
			break
		}
		assert.Equal(t, errFailed, err)
		assert.Nil(t, got)

		// Nothing is left of the bucket, and the old ones still work.
		assert.False(t, root.Get([]byte("new")).Found)
		assert.Equal(t, before, raw.Len())
		assert.Equal(t, []byte("1"), sub.Get([]byte("a")).Value)

		// It can be created afterwards.
		fs.putAt = 0
		got, err = root.NestedStore([]byte("new"))
		assert.NoError(t, err)
		assert.NoError(t, got.Put([]byte("k"), []byte("v")))
	}
}

func TestDeleteErrors(t *testing.T) {
	// Fail each Delete in turn until deleting the bucket works.
	for n := 1; ; n++ {
		fs := &failing{}
		root, _, _, _ := setup(t, fs)
		fs.deleteAt = n
		err := root.Delete([]byte("sub"))
		if err == nil {
			assert.Greater(t, n, 1)
			break
		}
		assert.Equal(t, errFailed, err)
	}

	// Saving the names of the buckets that are left can fail.
	fs := &failing{}
	root, _, _, _ := setup(t, fs)
	fs.putAt = 1
	assert.Equal(t, errFailed, root.Delete([]byte("sub")))

	// So can deleting a value.
	fs.deleteAt = 1
	assert.Equal(t, errFailed, root.Delete([]byte("keep")))
}
