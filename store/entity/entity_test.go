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

var foopath = StringToPath("base/foo")

func (*foo) EntPath() [][]byte {
	return foopath
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

	es := EntStore{
		Store:        s,
		Pather:       EntPathByEntPather{},
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
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
}

func TestGetSlice(t *testing.T) {
	s, err := ephemeral.Factory(bytebtree.New, 1).Store([]byte("test"))
	assert.NoError(t, err)

	es := EntStore{
		Store:        s,
		Pather:       EntPathByEntPather{},
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
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
		expected = append(expected, f)
	}

	var slc []*foo
	err = es.GetSlice(nil, &slc)
	assert.NoError(t, err)
	assert.Equal(t, expected, slc)
}
