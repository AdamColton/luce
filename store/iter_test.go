package store_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	s := lerr.Must(ephemeral.Factory(bytebtree.New, 1).FlatStore([]byte("test")))
	for i := byte(0); i < 10; i++ {
		j := i * i
		s.Put([]byte{i}, []byte{j})
	}

	i := store.NewIter(s)
	i.Filter = func(b []byte) bool {
		return b[0]%2 == 0
	}
	assert.Equal(t, -1, i.I)
	assert.Nil(t, i.Key)
	assert.False(t, i.Done())
	k, r, d := i.CurVal()
	assert.False(t, d)
	assert.Equal(t, []byte{0}, k)
	assert.Equal(t, []byte{0}, r.Value)
	assert.Equal(t, 0, i.Idx())

	i.Next()
	k, r, d = i.CurVal()
	assert.False(t, d)
	assert.Equal(t, []byte{2}, k)
	assert.Equal(t, []byte{4}, r.Value)

	timeout.After(5, func() {
		c := liter.ForIdx(i, func([]byte, int) {})
		assert.Equal(t, 5, i.Idx())
		assert.Equal(t, 4, c)
		k, r, d := i.CurVal()
		assert.True(t, d)
		assert.Nil(t, k)
		assert.Equal(t, store.Record{}, r)
	})
}
