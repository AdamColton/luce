package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

type kvs []lmap.KeyVal[string, int]

func (k kvs) Each(fn lmap.EachFunc[int, lmap.KeyVal[string, int]]) {
	done := false
	for i, kv := range k {
		fn(i, kv, &done)
		if done {
			return
		}
	}
}

func (k kvs) Len() int { return len(k) }

// noLen is an IdxEacher that is not a Lener.
type noLen struct {
	k kvs
}

func (n noLen) Each(fn lmap.EachFunc[int, lmap.KeyVal[string, int]]) {
	n.k.Each(fn)
}

func testKVs() kvs {
	return kvs{lmap.NewKV("a", 1), lmap.NewKV("b", 2), lmap.NewKV("a", 3)}
}

func TestEachType(t *testing.T) {
	calls := 0
	e := lmap.Each[string, int]{
		Func: func(fn lmap.EachFunc[string, int]) {
			calls++
			done := false
			fn("a", 1, &done)
		},
		L: 5,
	}
	var got map[string]int = map[string]int{}
	e.Each(func(k string, v int, done *bool) { got[k] = v })
	assert.Equal(t, map[string]int{"a": 1}, got)
	assert.Equal(t, 1, calls)
	assert.Equal(t, 5, e.Len())
}

func TestFromEach(t *testing.T) {
	m := lmap.FromEach(testKVs().Each)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, m.Map())
}

func TestFromEacher(t *testing.T) {
	expected := map[string]int{"a": 3, "b": 2}
	assert.Equal(t, expected, lmap.FromEacher[string, int](testKVs()).Map())
	assert.Equal(t, expected, lmap.FromEacher[string, int](noLen{testKVs()}).Map())
}

func TestAppendEach(t *testing.T) {
	m := lmap.New(map[string]int{"z": 26})
	got := m.AppendEacher(testKVs())
	assert.Equal(t, map[string]int{"a": 3, "b": 2, "z": 26}, m.Map())
	assert.Equal(t, m.Map(), got.Map())

	m = lmap.New(map[string]int{})
	m.AppendEach(testKVs().Each)
	assert.Equal(t, map[string]int{"a": 3, "b": 2}, m.Map())
}
