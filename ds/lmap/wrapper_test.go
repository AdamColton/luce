package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestWrapped(t *testing.T) {
	m := lmap.New(map[int]string{1: "1", 2: "2", 3: "3"})
	_, ok := m.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	w := lmap.Wrap(m)
	_, ok = w.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	assert.Equal(t, w.Mapper, w.Wrapped())
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

	assert.Equal(t, "apple", m.MustPop('a'))

	defer func() {
		assert.Equal(t, "failed to pop key: 97", recover().(error).Error())
	}()
	m.MustPop('a')
	t.Error("should not reach")
}

func TestAll(t *testing.T) {
	m := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	}
	lm := lmap.New(m)
	got := map[string]int{}
	lm.All(func(k string, v int) {
		got[k] = v
	})

	assert.Equal(t, m, got)
}
