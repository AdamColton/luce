package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

// pairs is an ordered Eacher so results are deterministic.
type pairs []struct {
	k string
	v int
}

func (p pairs) Each(fn slice.EachFn[string, int]) {
	done := false
	for _, kv := range p {
		fn(kv.k, kv.v, &done)
		if done {
			return
		}
	}
}

// lenPairs adds Len so it fulfills slice.Lener.
type lenPairs struct{ pairs }

func (lp lenPairs) Len() int { return len(lp.pairs) }

// wrapped fulfills upgrade.Wrapper but not slice.Lener, the Lener is only
// reachable through Wrapped.
type wrapped struct{ inner lenPairs }

func (w wrapped) Each(fn slice.EachFn[string, int]) { w.inner.Each(fn) }

func (w wrapped) Wrapped() any { return w.inner }

func testPairs() pairs {
	return pairs{{"a", 3}, {"b", 1}, {"c", 4}}
}

func TestEacherKey(t *testing.T) {
	keys := slice.EacherKey[string, int](testPairs())
	assert.Equal(t, slice.Slice[string]{"a", "b", "c"}, keys)

	keys = slice.EacherKey[string, int](lenPairs{testPairs()})
	assert.Equal(t, slice.Slice[string]{"a", "b", "c"}, keys)
	assert.Equal(t, 3, cap(keys))

	keys = slice.EacherKey[string, int](pairs{})
	assert.Empty(t, keys)
}

func TestEacherVal(t *testing.T) {
	vals := slice.EacherVal[string, int](testPairs())
	assert.Equal(t, slice.Slice[int]{3, 1, 4}, vals)

	vals = slice.EacherVal[string, int](wrapped{lenPairs{testPairs()}})
	assert.Equal(t, slice.Slice[int]{3, 1, 4}, vals)
	assert.Equal(t, 3, cap(vals))

	vals = slice.EacherVal[string, int](pairs{})
	assert.Empty(t, vals)
}
