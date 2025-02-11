package quicknested_test

import (
	"testing"

	"github.com/adamcolton/luce/store/ephemeral/quicknested"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// The buffer size is only where the store starts, so any size works.
	for _, size := range []int{0, 1, 100} {
		testsuite.TestAll(t, quicknested.New(size))
	}
}

func TestNested(t *testing.T) {
	fac := quicknested.New(1)
	users, err := fac.NestedStore([]byte("users"))
	assert.NoError(t, err)
	assert.NoError(t, users.Put([]byte("ada"), []byte("Ada Lovelace")))

	// A store can hold other stores, as many levels deep as needed.
	prefs, err := users.NestedStore([]byte("prefs"))
	assert.NoError(t, err)
	assert.NoError(t, prefs.Put([]byte("theme"), []byte("dark")))

	// Asking for a store again gets what was put in it.
	again, err := fac.NestedStore([]byte("users"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("Ada Lovelace"), again.Get([]byte("ada")).Value)
	assert.Equal(t, []byte("dark"), again.Get([]byte("prefs")).Store.Get([]byte("theme")).Value)
	assert.Equal(t, 2, again.Len())

	// Stores with different names are separate.
	other, err := fac.NestedStore([]byte("groups"))
	assert.NoError(t, err)
	assert.False(t, other.Get([]byte("ada")).Found)
}
