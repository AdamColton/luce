package store_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/stretchr/testify/assert"
)

func TestGetStores(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 10)
	s, err := store.GetStoresStr(f, "this", "is", "a", "test")
	assert.NoError(t, err)
	assert.Len(t, s, 4)
}
