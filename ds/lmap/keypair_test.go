package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestKeyPairLess(t *testing.T) {
	a := lmap.KeyPair[string, int]{Key: "a", Val: 2}
	b := lmap.KeyPair[string, int]{Key: "b", Val: 1}
	keyLess := func(i, j string) bool { return i < j }
	valLess := func(i, j int) bool { return i < j }

	assert.True(t, lmap.KeyLessKP[int](keyLess)(a, b))
	assert.False(t, lmap.KeyLessKP[int](keyLess)(b, a))
	assert.False(t, lmap.ValLessKP[string](valLess)(a, b))
	assert.True(t, lmap.ValLessKP[string](valLess)(b, a))

	m := lmap.New(map[string]int{})
	assert.True(t, m.KeyLessKP(keyLess)(a, b))
	assert.False(t, m.ValLessKP(valLess)(a, b))
}

func TestWrapperSlice(t *testing.T) {
	m := lmap.New(map[string]int{"b": 1, "c": 3, "a": 2})

	got := m.Slice(m.KeyLessKP(func(i, j string) bool { return i < j }), nil)
	expected := []lmap.KeyPair[string, int]{{Key: "a", Val: 2}, {Key: "b", Val: 1}, {Key: "c", Val: 3}}
	assert.Equal(t, expected, []lmap.KeyPair[string, int](got))

	got = m.Slice(m.ValLessKP(func(i, j int) bool { return i < j }), nil)
	expected = []lmap.KeyPair[string, int]{{Key: "b", Val: 1}, {Key: "a", Val: 2}, {Key: "c", Val: 3}}
	assert.Equal(t, expected, []lmap.KeyPair[string, int](got))

	// Without less the order is not defined, but every pair is there.
	buf := make([]lmap.KeyPair[string, int], 0, 3)
	got = m.Slice(nil, buf)
	assert.ElementsMatch(t, expected, []lmap.KeyPair[string, int](got))
	assert.Equal(t, 3, cap(got))
}
