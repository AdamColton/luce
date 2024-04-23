package store_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

func TestGetStores(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 10)
	s, err := store.GetStoresStr(f, "this", "is", "a", "test")
	assert.NoError(t, err)
	assert.Equal(t, 4, s.Len())
}

func TestSlice(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 10)
	s, err := f.Store([]byte("test"))
	assert.NoError(t, err)

	expected := slice.Slice[[]byte]{[]byte("this"), []byte("is"), []byte("a"), []byte("test")}
	for _, k := range expected {
		s.Put(k, []byte{0})
	}
	expected.Sort(func(i, j []byte) bool {
		return bytes.Compare(i, j) == -1
	})

	got := store.Slice(s)
	assert.Equal(t, expected, got)
}
