package entity

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

type foo struct {
	key string
	S   string
	F   float64
	I   int
}

func (f *foo) EntKey() []byte {
	return []byte(f.key)
}
func (f *foo) SetEntKey(key []byte) {
	f.key = string(key)
}

func TestRoundTrip(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).Store([]byte("test"))
	assert.NoError(t, err)

	es := EntStore[*foo]{
		Store:        s,
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
		Init: func() *foo {
			return &foo{}
		},
	}

	f := &foo{
		key: "key1",
		S:   "just a string",
		F:   3.1415,
		I:   42,
	}
	_, err = es.Put(f, nil)
	assert.NoError(t, err)

	f2 := &foo{
		key: "key1",
	}
	err = es.Load(f2)
	assert.NoError(t, err)

	assert.Equal(t, f, f2)
	found, f3, err := es.Get(f.EntKey())
	assert.True(t, found)
	assert.NoError(t, err)
	assert.Equal(t, f, f3)

	f.key = "not in store"
	found, f3, err = es.Get(f.EntKey())
	assert.False(t, found)
	assert.Nil(t, f3)
	assert.NoError(t, err)
	err = es.Load(f)
	assert.Equal(t, ErrKeyNotFound, err)

	f.key = "corrupted data"
	es.Store.Put(f.EntKey(), []byte("this is not json"))
	found, f3, err = es.Get(f.EntKey())
	assert.True(t, found)
	assert.Equal(t, es.Init(), f3)
	assert.Error(t, err)
	err = es.Load(f)
	assert.Error(t, err)
	assert.True(t, err != ErrKeyNotFound)
}

func TestGetSlice(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).Store([]byte("test"))
	assert.NoError(t, err)

	es := EntStore[*foo]{
		Store:        s,
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
		Init: func() *foo {
			return &foo{}
		},
	}

	var expected []*foo
	for i := 1; i < 20; i++ {
		f := &foo{
			key: fmt.Sprintf("key-%02d", i),
			S:   fmt.Sprintf("I am Foo #%d", i),
			F:   float64(i) * 3.1415,
			I:   i,
		}
		_, err = es.Put(f, nil)
		assert.NoError(t, err)
		if i != 10 {
			expected = append(expected, f)
		}
	}

	kf := func(b []byte) bool {
		return string(b) != "key-10"
	}

	slc, err := es.GetSlice(kf, nil)
	assert.NoError(t, err)
	assert.Equal(t, expected, slc)

	es.Store.Put(expected[0].EntKey(), []byte("this is not json"))
	slc, err = es.GetSlice(kf, nil)
	assert.Error(t, err)
	assert.Nil(t, slc)
}

func TestIndex(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 1)
	s, err := f.Store([]byte("ents"))
	assert.NoError(t, err)
	i, err := f.Store([]byte("idxs"))
	assert.NoError(t, err)

	es := EntStore[*foo]{
		Store:        s,
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
		Init: func() *foo {
			return &foo{}
		},
		IdxStore: i,
	}
	es.AddIndex(BaseIndexer[*foo]{
		IndexName: "byS",
		Fn: func(f *foo) []byte {
			return []byte(f.S)
		},
	})
	es.AddIndex(BaseIndexer[*foo]{
		IndexName: "mod2",
		Fn: func(f *foo) []byte {
			return []byte{byte(f.I % 2)}
		},
		M: true,
	})

	for i := 1; i < 20; i++ {
		f := &foo{
			key: fmt.Sprintf("key-%02d", i),
			S:   fmt.Sprintf("I am Foo #%d", i),
			F:   float64(i) * 3.1415,
			I:   i,
		}
		_, err = es.Put(f, nil)
		assert.NoError(t, err)
	}

	ents, err := es.Index("mod2", []byte{1})
	assert.NoError(t, err)
	assert.Len(t, ents, 10)

	expected := "I am Foo #3"
	ents, err = es.Index("byS", []byte(expected))
	assert.NoError(t, err)
	assert.Len(t, ents, 1)
	assert.Equal(t, expected, ents[0].S)

	e := ents[0]
	e.S = "FOOOOO"
	e.I = 100
	es.Put(e, nil)

	ents, err = es.Index("mod2", []byte{1})
	assert.NoError(t, err)
	assert.Len(t, ents, 9)

	ents, err = es.Index("byS", []byte(e.S))
	assert.NoError(t, err)
	assert.Len(t, ents, 1)
	assert.Equal(t, e, ents[0])

}
