// Package testsuite provides shared tests for implementations of the store
// interfaces.
package testsuite

import (
	"bytes"
	"math/rand"
	"sort"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/stretchr/testify/assert"
)

// TestAll runs every test in the package against factory. A FlatStore is the
// NestedStore of the same name, as store.Flattener does it. The tests use random
// bytes as keys, so an implementation that turns keys into file names or other
// text has to encode them.
func TestAll(t *testing.T, factory store.NestedFactory) {
	TestFlat(t, store.Flattener{factory})
	TestBuckets(t, factory)
	TestBucketDataCollision(t, factory)
	TestIteration(t, store.Flattener{factory})
}

// TestFlat tests Put, Get, Next, Delete and Len on a FlatStore, and that asking
// the factory for the same name again gives a store with the same records.
func TestFlat(t *testing.T, factory store.FlatFactory) {
	s, err := factory.FlatStore([]byte("TestBasic"))
	assert.NoError(t, err)

	k, v := []byte{1, 2, 3}, []byte{4, 5, 6, 7, 8, 9, 10}
	err = s.Put(k, v)
	assert.NoError(t, err)
	assert.Equal(t, v, s.Get(k).Value)
	assert.Equal(t, 1, s.Len())

	k2 := []byte{3, 2, 1}
	v2 := []byte{10, 9, 8, 7, 6, 5, 4}
	s.Put(k2, v2)
	assert.Equal(t, k2, s.Next(k))
	assert.Equal(t, 2, s.Len())

	s.Delete(k)
	assert.Nil(t, s.Get(k).Value)
	assert.Equal(t, v2, s.Get(k2).Value)
	s2, err := factory.FlatStore([]byte("TestBasic"))
	assert.NoError(t, err)
	assert.Equal(t, v2, s2.Get(k2).Value)

	assert.Equal(t, 1, s.Len())
}

// TestBuckets tests Sub-Stores: they can be created, found again through the
// parent by NestedStore and by Get, and deleted. It deletes the bucket it makes.
func TestBuckets(t *testing.T, factory store.NestedFactory) {
	s, err := factory.NestedStore([]byte("TestBuckets"))
	assert.NoError(t, err)

	k := []byte("bucketName")
	bkt := make([]byte, 10)
	_, err = rand.Read(bkt)
	assert.NoError(t, err)

	assert.NoError(t, s.Put(k, bkt))
	s2, err := s.NestedStore(bkt)
	assert.NoError(t, err)

	k2 := []byte("foo")
	v2 := []byte("bar")
	assert.NoError(t, s2.Put(k2, v2))

	assert.Equal(t, bkt, s.Get(k).Value)

	s3, err := s.NestedStore(bkt)
	assert.NoError(t, err)
	assert.Equal(t, v2, s3.Get(k2).Value)

	r := s.Get(bkt)
	assert.True(t, r.Found)
	assert.NotNil(t, r.Store)
	assert.Equal(t, v2, r.Store.Get(k2).Value)

	s.Delete(bkt)
	r = s.Get(bkt)
	assert.False(t, r.Found)
	assert.Nil(t, r.Store)

}

// TestBucketDataCollision tests that a key holds either a value or a Sub-Store:
// creating a Sub-Store at the key of a value is an error, and so is putting a
// value at the key of a Sub-Store.
func TestBucketDataCollision(t *testing.T, factory store.NestedFactory) {
	s, err := factory.NestedStore([]byte("TestBucketDataCollision"))
	assert.NoError(t, err)

	k := make([]byte, 10)
	_, err = rand.Read(k)
	assert.NoError(t, err)
	err = s.Put(k, []byte("foo"))
	assert.NoError(t, err)
	_, err = s.NestedStore(k)
	assert.Error(t, err)

	k = make([]byte, 10)
	_, err = rand.Read(k)
	assert.NoError(t, err)
	_, err = s.NestedStore(k)
	assert.NoError(t, err)
	err = s.Put(k, []byte("foo"))
	assert.Error(t, err)
}

// TestIteration tests that Next walks the keys of a FlatStore in byte order,
// starting from nil and ending with nil.
func TestIteration(t *testing.T, factory store.FlatFactory) {
	s, err := factory.FlatStore([]byte("TestIteration"))
	assert.NoError(t, err)

	vals := [][]byte{
		{3, 4, 5},
		{1, 1, 0},
		{2, 0, 0},
		{1, 1, 1},
		{1, 0, 0},
	}
	for _, v := range vals {
		err = s.Put(v, v)
		assert.NoError(t, err)
	}

	// == projects.Code.luce.store ==
	// [ ] TestIteration should include Sub-Stores
	//  Next returns the keys of Sub-Stores along with the keys of values, but
	//  this test takes a FlatFactory, so it can't make one. It needs a version
	//  that takes a NestedFactory.

	sort.Slice(vals, func(i, j int) bool {
		return bytes.Compare(vals[i], vals[j]) == -1
	})

	// It stops after one key too many, so a Next that never returns nil fails
	// the test instead of running forever.
	var got [][]byte
	for cur := s.Next(nil); cur != nil && len(got) <= len(vals); cur = s.Next(cur) {
		got = append(got, cur)
	}
	assert.Equal(t, vals, got)
}
