package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
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

func TestKeys(t *testing.T) {
	m := lmap.New(map[rune]string{
		'a': "apple",
		'b': "banana",
		'c': "cantaloupe",
	})
	ks := m.Keys(nil).Sort(slice.LT[rune]())
	expected := slice.Slice[rune]{'a', 'b', 'c'}
	assert.Equal(t, expected, ks)
}

func TestVals(t *testing.T) {
	m := lmap.New(map[rune]string{
		'a': "apple",
		'b': "banana",
		'c': "cantaloupe",
	})
	ks := m.Vals(nil).Sort(slice.LT[string]())
	expected := slice.Slice[string]{"apple", "banana", "cantaloupe"}
	assert.Equal(t, expected, ks)
}

func TestSortKeys(t *testing.T) {
	m := map[rune]string{
		'a': "apple",
		'b': "banana",
		'c': "cantaloupe",
	}
	ks := lmap.SortKeys(m)
	expected := slice.Slice[rune]{'a', 'b', 'c'}
	assert.Equal(t, expected, ks)
}

func TestWrapNew(t *testing.T) {
	m := lmap.Empty[int, string](0)
	_, ok := m.WrapNew().Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	s := lmap.EmptySafe[int, string](0)
	_, ok = s.WrapNew().Mapper.(*lmap.Safe[int, string])
	assert.True(t, ok)
}
