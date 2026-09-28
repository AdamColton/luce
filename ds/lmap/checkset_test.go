package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestCheckSet(t *testing.T) {
	m := lmap.New(map[string]int{})

	assert.NoError(t, m.CheckSet("a", 1))
	assert.Equal(t, 1, m.GetVal("a"))

	assert.Equal(t, lmap.ErrCollision, m.CheckSet("a", 2))
	assert.Equal(t, 1, m.GetVal("a"))

	assert.Equal(t, lmap.ErrCollision, lmap.CheckSet[string, int](m, "a", 3))
	assert.NoError(t, lmap.CheckSet[string, int](m, "b", 4))
	assert.Equal(t, 4, m.GetVal("b"))
}
