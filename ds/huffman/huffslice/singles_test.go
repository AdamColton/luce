package huffslice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/huffman/huffslice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func decode(s *huffslice.Slice[string]) []string {
	var out []string
	liter.For(s.Iter(), func(v string) { out = append(out, v) })
	return out
}

func encode(token string, vals ...string) *huffslice.Slice[string] {
	e := huffslice.NewEncoder(len(vals), token)
	e.Slice = append(e.Slice, vals...)
	return e.Encode()
}

func TestSinglesRoundTrip(t *testing.T) {
	// The first and last values occur once, so they are singles.
	vals := []string{"a", "b", "b", "c", "b", "d"}
	s := encode("", vals...)
	assert.Equal(t, []string{"a", "c", "d"}, []string(s.Singles))
	assert.Equal(t, vals, decode(s))

	// The Cur of the iterator is a single at the start.
	it := s.Iter()
	v, done := it.Cur()
	assert.Equal(t, "a", v)
	assert.False(t, done)
	v, done = it.Next()
	assert.Equal(t, "b", v)
	assert.False(t, done)
}

func TestSingleTokenInData(t *testing.T) {
	// A value equal to SingleToken is a single even if it occurs many times.
	vals := []string{"x", "y", "x", "y", "x"}
	s := encode("x", vals...)
	assert.Equal(t, []string{"x", "x", "x"}, []string(s.Singles))
	assert.Equal(t, vals, decode(s))
}

func TestAllSingles(t *testing.T) {
	// No value repeats, so Encoded is empty and Iter reads the Singles.
	vals := []string{"a", "b", "c"}
	s := encode("", vals...)
	assert.Zero(t, s.Encoded.Ln)
	assert.Equal(t, vals, decode(s))

	assert.Empty(t, decode(encode("")))
}
