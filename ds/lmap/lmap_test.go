package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	m := lmap.Map[int, string]{1: "1", 2: "2", 3: "3"}
	assert.Equal(t, len(m), m.Len())
}

func TestPop(t *testing.T) {
	m := lmap.Map[rune, string]{
		'a': "apple",
		'b': "banana",
		'c': "cantaloupe",
	}
	got, ok := m.Pop('b')
	assert.True(t, ok)
	assert.Equal(t, "banana", got)
	_, ok = m['b']
	assert.False(t, ok)

	got, ok = m.Pop('d')
	assert.False(t, ok)
	assert.Equal(t, "", got)
}
