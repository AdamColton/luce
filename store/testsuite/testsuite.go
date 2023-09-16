package testsuite

import (
	"bytes"
	"math/rand"
	"sort"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T, factory store.Factory) {
	TestBasic(t, factory)
	TestBuckets(t, factory)
	TestBucketDataCollision(t, factory)
	TestIteration(t, factory)
}

func TestBasic(t *testing.T, factory store.Factory) {
	s, err := factory.Store([]byte("TestBasic"))
	assert.NoError(t, err)

	k, v := []byte{1, 2, 3}, []byte{4, 5, 6, 7, 8, 9, 10}
	err = s.Put(k, v)
	assert.NoError(t, err)
	assert.Equal(t, v, s.Get(k).Value)

	k2 := []byte{3, 2, 1}
	v2 := []byte{10, 9, 8, 7, 6, 5, 4}
	s.Put(k2, v2)
	assert.Equal(t, k2, s.Next(k))

	s.Delete(k)
	assert.Nil(t, s.Get(k).Value)
	assert.Equal(t, v2, s.Get(k2).Value)
	s2, err := factory.Store([]byte("TestBasic"))
	assert.NoError(t, err)
	assert.Equal(t, v2, s2.Get(k2).Value)

}

func TestBuckets(t *testing.T, factory store.Factory) {
	s, err := factory.Store([]byte("TestBuckets"))
	assert.NoError(t, err)

	k := []byte("bucketName")
	bkt := make([]byte, 10)
	_, err = rand.Read(bkt)
	assert.NoError(t, err)

	assert.NoError(t, s.Put(k, bkt))
	s2, err := s.Store(bkt)
	assert.NoError(t, err)

	k2 := []byte("foo")
	v2 := []byte("bar")
	assert.NoError(t, s2.Put(k2, v2))

	assert.Equal(t, bkt, s.Get(k).Value)

	s3, err := s.Store(bkt)
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

func TestBucketDataCollision(t *testing.T, factory store.Factory) {
	s, err := factory.Store([]byte("TestBucketDataCollision"))
	assert.NoError(t, err)

	k := make([]byte, 10)
	_, err = rand.Read(k)
	assert.NoError(t, err)
	err = s.Put(k, []byte("foo"))
	assert.NoError(t, err)
	_, err = s.Store(k)
	assert.Error(t, err)

	k = make([]byte, 10)
	_, err = rand.Read(k)
	assert.NoError(t, err)
	_, err = s.Store(k)
	assert.NoError(t, err)
	err = s.Put(k, []byte("foo"))
	assert.Error(t, err)
}

func TestIteration(t *testing.T, factory store.Factory) {
	s, err := factory.Store([]byte("TestIteration"))
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

	bkt := []byte{2, 2, 2}
	vals = append(vals, bkt)
	_, err = s.Store(bkt)
	assert.NoError(t, err)

	sort.Slice(vals, func(i, j int) bool {
		return bytes.Compare(vals[i], vals[j]) == -1
	})

	i := 0
	for cur := s.Next(nil); cur != nil; cur = s.Next(cur) {
		if !assert.NotEqual(t, i, len(vals)) {
			break
		}
		assert.Equal(t, cur, vals[i])
		i++
	}
	assert.Equal(t, i, len(vals))
}
