package store_test

import (
	"errors"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/stretchr/testify/assert"
)

// nested is a NestedStore that only remembers its name. Its factory fails for
// the name "bad".
type nested struct {
	name string
}

func (nested) Put(key, value []byte) error { return nil }
func (nested) Get(key []byte) store.Record { return store.Record{} }
func (nested) Next(key []byte) []byte      { return nil }
func (nested) Delete(key []byte) error     { return nil }
func (nested) Len() int                    { return 0 }

func (n nested) NestedStore(name []byte) (store.NestedStore, error) {
	if string(name) == "bad" {
		return nil, errors.New("bad name")
	}
	return nested{name: string(name)}, nil
}

func TestFlattener(t *testing.T) {
	var f store.FlatFactory = store.Flattener{nested{}}

	s, err := f.FlatStore([]byte("good"))
	assert.NoError(t, err)
	assert.Equal(t, nested{name: "good"}, s)

	s, err = f.FlatStore([]byte("bad"))
	assert.EqualError(t, err, "bad name")
	assert.Nil(t, s)
}
