package ephemeral_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/ds/idx/byteid/bytemap"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

var indexes = map[string]byteid.IndexFactory{
	"bytebtree": bytebtree.New,
	"bytemap":   bytemap.New,
}

func TestAll(t *testing.T) {
	// bufferSize is only the starting size, so it works too small, and too large.
	for name, idx := range indexes {
		for _, size := range []int{0, 1, 100} {
			t.Run(fmt.Sprintf("%s/%d", name, size), func(t *testing.T) {
				f := ephemeral.Factory(idx, size)
				testsuite.TestFlat(t, f)
				testsuite.TestIteration(t, f)
			})
		}
	}
}

func TestFactory(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 2)

	// More stores than the bufferSize; each name is a separate store.
	names := []string{"a", "b", "c", "d"}
	for _, name := range names {
		s, err := f.FlatStore([]byte(name))
		assert.NoError(t, err)
		assert.NoError(t, s.Put([]byte("name"), []byte(name)))
	}
	for _, name := range names {
		s, err := f.FlatStore([]byte(name))
		assert.NoError(t, err)
		assert.Equal(t, []byte(name), s.Get([]byte("name")).Value)
		assert.Equal(t, 1, s.Len())
	}

	// The name is copied, so the caller can reuse the buffer.
	name := []byte("reused")
	s, _ := f.FlatStore(name)
	s.Put([]byte("k"), []byte("v"))
	copy(name, "xxxxxx")
	s, _ = f.FlatStore([]byte("reused"))
	assert.Equal(t, []byte("v"), s.Get([]byte("k")).Value)
}

func TestReuse(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).FlatStore([]byte("reuse"))
	assert.NoError(t, err)

	// Put copies the key and the value.
	key, val := []byte("key"), []byte("value")
	assert.NoError(t, s.Put(key, val))
	copy(key, "xxx")
	copy(val, "xxxxx")
	assert.Equal(t, store.Record{Found: true, Value: []byte("value")}, s.Get([]byte("key")))
	assert.False(t, s.Get([]byte("xxx")).Found)

	// Get and Next return copies.
	s.Get([]byte("key")).Value[0] = 'x'
	assert.Equal(t, []byte("value"), s.Get([]byte("key")).Value)
	s.Next(nil)[0] = 'x'
	assert.Equal(t, []byte("key"), s.Next(nil))
}

func TestReplaceDelete(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).FlatStore([]byte("replace"))
	assert.NoError(t, err)

	assert.NoError(t, s.Put([]byte("a"), []byte("1")))
	assert.NoError(t, s.Put([]byte("a"), []byte("2")))
	assert.Equal(t, 1, s.Len())
	assert.Equal(t, []byte("2"), s.Get([]byte("a")).Value)

	// A key that holds an empty value is still found.
	assert.NoError(t, s.Put([]byte("empty"), nil))
	assert.Equal(t, store.Record{Found: true}, s.Get([]byte("empty")))

	assert.NoError(t, s.Delete([]byte("a")))
	assert.NoError(t, s.Delete([]byte("a")))
	assert.False(t, s.Get([]byte("a")).Found)
	assert.Equal(t, 1, s.Len())

	// The slot is reused.
	assert.NoError(t, s.Put([]byte("b"), []byte("3")))
	assert.Equal(t, []byte("3"), s.Get([]byte("b")).Value)
	assert.Equal(t, []byte("b"), s.Next([]byte("a")))
	assert.Equal(t, []byte("empty"), s.Next([]byte("b")))
	assert.Nil(t, s.Next([]byte("empty")))
}
