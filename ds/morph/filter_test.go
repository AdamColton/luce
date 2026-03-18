package morph_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

var (
	keyIsA  = filter.Filter[string](func(s string) bool { return s == "a" })
	valBig  = filter.Filter[int](func(i int) bool { return i > 1 })
	outLong = filter.Filter[string](func(s string) bool { return len(s) > 2 })
)

type kvResult struct {
	out     string
	include bool
}

func runKV(f morph.KeyVal[string, int, string], k string, v int) kvResult {
	out, include := f(k, v)
	return kvResult{out, include}
}

func TestKeyValAllFilters(t *testing.T) {
	calls := 0
	kva := morph.NewKeyValAll(func(k string, v int) string {
		calls++
		return k + strconv.Itoa(v)
	})

	assert.Equal(t, kvResult{"a1", true}, runKV(kva.FilterKey(keyIsA), "a", 1))
	assert.False(t, runKV(kva.FilterKey(keyIsA), "b", 1).include)
	assert.Equal(t, kvResult{"a2", true}, runKV(kva.FilterVal(valBig), "a", 2))
	assert.False(t, runKV(kva.FilterVal(valBig), "a", 1).include)
	assert.Equal(t, kvResult{"a12", true}, runKV(kva.FilterOut(outLong), "a", 12))
	assert.False(t, runKV(kva.FilterOut(outLong), "a", 1).include)

	// FilterKey and FilterVal only call the KeyValAll for included pairs.
	calls = 0
	runKV(kva.FilterKey(keyIsA), "b", 1)
	runKV(kva.FilterVal(valBig), "a", 1)
	assert.Zero(t, calls)
}

func TestKeyValFilters(t *testing.T) {
	// excludes any value of 5
	kv := morph.NewKeyVal(func(k string, v int) (string, bool) { return k + strconv.Itoa(v), v != 5 })

	assert.Equal(t, kvResult{"a1", true}, runKV(kv.FilterKey(keyIsA), "a", 1))
	assert.False(t, runKV(kv.FilterKey(keyIsA), "b", 1).include)
	assert.False(t, runKV(kv.FilterKey(keyIsA), "a", 5).include)

	assert.Equal(t, kvResult{"a2", true}, runKV(kv.FilterVal(valBig), "a", 2))
	assert.False(t, runKV(kv.FilterVal(valBig), "a", 1).include)
	assert.False(t, runKV(kv.FilterVal(valBig), "a", 5).include)

	assert.Equal(t, kvResult{"a12", true}, runKV(kv.FilterOut(outLong), "a", 12))
	assert.False(t, runKV(kv.FilterOut(outLong), "a", 1).include)
	assert.False(t, runKV(kv.FilterOut(outLong), "a", 5).include)
}

func TestValFilters(t *testing.T) {
	va := morph.NewValAll(func(v int) string { return "v" + strconv.Itoa(v) })

	out, include := va.FilterVal(valBig)(2)
	assert.Equal(t, "v2", out)
	assert.True(t, include)
	_, include = va.FilterVal(valBig)(1)
	assert.False(t, include)

	out, include = va.FilterOut(outLong)(12)
	assert.Equal(t, "v12", out)
	assert.True(t, include)
	_, include = va.FilterOut(outLong)(1)
	assert.False(t, include)

	// excludes any value of 5
	v := morph.NewVal(func(v int) (string, bool) { return "v" + strconv.Itoa(v), v != 5 })
	out, include = v.FilterVal(valBig)(2)
	assert.Equal(t, "v2", out)
	assert.True(t, include)
	_, include = v.FilterVal(valBig)(1)
	assert.False(t, include)
	_, include = v.FilterVal(valBig)(5)
	assert.False(t, include)

	out, include = v.FilterOut(outLong)(12)
	assert.Equal(t, "v12", out)
	assert.True(t, include)
	_, include = v.FilterOut(outLong)(1)
	assert.False(t, include)
	_, include = v.FilterOut(outLong)(5)
	assert.False(t, include)
}
