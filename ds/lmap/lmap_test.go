package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	m := lmap.New(map[int]string{1: "1", 2: "2", 3: "3"})

	assert.Equal(t, 3, m.Len())

	v, found := m.Get(1)
	assert.Equal(t, "1", v)
	assert.Equal(t, "1", m.GetVal(1))
	assert.True(t, found)

	v, found = m.Get(4)
	assert.Equal(t, "", v)
	assert.Equal(t, "", m.GetVal(4))
	assert.False(t, found)

	m.Set(4, "4")
	v, found = m.Get(4)
	assert.Equal(t, "4", v)
	assert.True(t, found)

	assert.Equal(t, 4, m.Len())

	m.Delete(1)
	v, found = m.Get(1)
	assert.Equal(t, "", v)
	assert.False(t, found)

	got := make(map[int]string)
	m.Each(func(key int, val string) (done bool) {
		got[key] = val
		return false
	})
	assert.Equal(t, m, lmap.New(got))
}

func TestPop(t *testing.T) {
	m := lmap.New(map[rune]string{
		'a': "apple",
		'b': "banana",
		'c': "cantaloupe",
	})
	got, ok := m.Pop('b')
	assert.True(t, ok)
	assert.Equal(t, "banana", got)
	_, ok = m.Get('b')
	assert.False(t, ok)

	got, ok = m.Pop('d')
	assert.False(t, ok)
	assert.Equal(t, "", got)
}
