package morph_test

import (
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/morph"
	"github.com/stretchr/testify/assert"
)

func testMap() lmap.Wrapper[string, int] {
	return lmap.New(map[string]int{"a": 1, "b": 2, "c": 3})
}

// noLen is an Eacher that is not a Lener.
type noLen struct {
	m lmap.Wrapper[string, int]
}

func (n noLen) Each(fn morph.EachFn[string, int]) {
	n.m.Each(fn)
}

func eachSorted(e morph.Eacher[int, string]) (idxs []int, vals []string) {
	e.Each(func(idx int, s string, done *bool) {
		idxs = append(idxs, idx)
		vals = append(vals, s)
	})
	sort.Ints(idxs)
	sort.Strings(vals)
	return
}

func TestKeyValAllEacher(t *testing.T) {
	kva := morph.NewKeyValAll(func(k string, v int) string { return k + strconv.Itoa(v) })
	e := kva.Eacher(testMap())
	idxs, vals := eachSorted(e)
	assert.Equal(t, []int{0, 1, 2}, idxs)
	assert.Equal(t, []string{"a1", "b2", "c3"}, vals)
	assert.Equal(t, 3, e.(morph.Lener).Len())

	// Setting done stops the iteration.
	calls := 0
	e.Each(func(idx int, s string, done *bool) {
		calls++
		*done = true
	})
	assert.Equal(t, 1, calls)

	// Without a Len, the Len is 0.
	e = kva.Eacher(noLen{testMap()})
	assert.Equal(t, 0, e.(morph.Lener).Len())
	_, vals = eachSorted(e)
	assert.Equal(t, []string{"a1", "b2", "c3"}, vals)
}

func TestKeyValEacher(t *testing.T) {
	kv := morph.NewKeyVal(func(k string, v int) (string, bool) { return k + strconv.Itoa(v), v != 2 })
	e := kv.Eacher(testMap())
	idxs, vals := eachSorted(e)
	// The index counts the included results.
	assert.Equal(t, []int{0, 1}, idxs)
	assert.Equal(t, []string{"a1", "c3"}, vals)
	// The Len is that of the source, not of the included results.
	assert.Equal(t, 3, e.(morph.Lener).Len())

	calls := 0
	e.Each(func(idx int, s string, done *bool) {
		calls++
		*done = true
	})
	assert.Equal(t, 1, calls)
}

func TestOnValOnKey(t *testing.T) {
	m := testMap()

	onVal := morph.OnVal[string](morph.NewValAll(func(v int) string { return strconv.Itoa(v * 10) }))
	got := lmap.FromEacher(onVal.Eacher(m))
	assert.Equal(t, map[string]string{"a": "10", "b": "20", "c": "30"}, got.Map())

	onKey := morph.OnKey[int](morph.NewValAll(strings.ToUpper))
	got2 := lmap.FromEacher(onKey.Eacher(m))
	assert.Equal(t, map[string]int{"A": 1, "B": 2, "C": 3}, got2.Map())
}
