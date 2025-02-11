package flatwrap_test

import (
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
