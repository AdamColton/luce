package roaring_test

import (
	"math/rand/v2"
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestSetAnd(t *testing.T) {
	// Which blocked addresses tried to connect today?
	blocked := roaring.New(ip(10, 0, 0, 1), ip(192, 168, 1, 7), ip(203, 0, 113, 9))
	today := roaring.New(ip(192, 168, 1, 7), ip(8, 8, 8, 8), ip(203, 0, 113, 9), ip(10, 0, 0, 2))

	caught := blocked.And(today, nil)
	assert.Equal(t, []uint32{ip(192, 168, 1, 7), ip(203, 0, 113, 9)}, []uint32(caught.Slice(nil)))
	assert.Equal(t, 3, blocked.Len(), "inputs are unchanged")
	assert.Equal(t, 4, today.Len())

	// Doing it again tomorrow, the result from today is reused as the buffer.
	tomorrow := roaring.New(ip(10, 0, 0, 1))
	caught2 := blocked.And(tomorrow, caught)
	assert.Same(t, caught, caught2)
	assert.Equal(t, []uint32{ip(10, 0, 0, 1)}, []uint32(caught2.Slice(nil)))
}

// fill adds every value in [lo, hi) with the given step, all in the /16
// network numbered key.
func fill(s *roaring.Set, key, lo, hi, step uint32) *roaring.Set {
	for v := lo; v < hi; v += step {
		s.Add(key<<16 | v)
	}
	return s
}

func TestSetAndPairings(t *testing.T) {
	// Container 0 is Array in both, 1 is Array∧Bitmap, 2 is Bitmap∧Array and
	// 3 is Bitmap∧Bitmap with a large overlap. In 4 two Bitmaps overlap in
	// only a few values, so the result goes back to an Array. 5 has no
	// overlap at all, and 6 and 7 are each only in one set.
	a := roaring.New()
	fill(a, 0, 0, 1000, 2)
	fill(a, 1, 0, 1000, 3)
	fill(a, 2, 0, 60000, 7)
	fill(a, 3, 0, 60000, 2)
	fill(a, 4, 0, 30000, 1)
	fill(a, 5, 0, 100, 2)
	fill(a, 6, 0, 10, 1)

	b := roaring.New()
	fill(b, 0, 0, 1000, 3)
	fill(b, 1, 0, 60000, 2)
	fill(b, 2, 0, 1000, 5)
	fill(b, 3, 0, 60000, 3)
	fill(b, 4, 29990, 60000, 1)
	fill(b, 5, 1, 100, 2)
	fill(b, 7, 0, 10, 1)

	got := a.And(b, nil)
	assert.Equal(t, andLset(a, b), []uint32(got.Slice(nil)))
	assert.Equal(t, got.Len(), len(got.Slice(nil)))
	assert.True(t, got.Contains(4<<16|29995))
	assert.False(t, got.Contains(5<<16|1))
	assert.Equal(t, got.Slice(nil), b.And(a, nil).Slice(nil), "And is symmetric")
}

func TestSetAndBuffer(t *testing.T) {
	// A buffer full of unrelated containers of both kinds is cleared and its
	// containers are reused for the result.
	buf := roaring.New()
	fill(buf, 9, 0, 60000, 1)
	fill(buf, 3, 0, 60000, 1)
	fill(buf, 0, 0, 50, 1)

	a := fill(fill(roaring.New(), 3, 0, 60000, 2), 0, 0, 1000, 1)
	b := fill(fill(roaring.New(), 3, 0, 60000, 3), 0, 500, 2000, 1)
	got := a.And(b, buf)
	assert.Same(t, buf, got)
	assert.Equal(t, andLset(a, b), []uint32(got.Slice(nil)))
	assert.False(t, got.Contains(9<<16), "old contents are gone")
}

func TestSetAndEmpty(t *testing.T) {
	var empty roaring.Set
	a := roaring.New(1, 2, 3)
	assert.Equal(t, 0, a.And(&empty, nil).Len())
	assert.Equal(t, 0, empty.And(a, nil).Len())
}

func TestSetAndRandom(t *testing.T) {
	rnd := rand.New(rand.NewPCG(3, 4))
	var buf *roaring.Set
	for range 20 {
		a, b := randomSet(rnd), randomSet(rnd)
		buf = a.And(b, buf)
		assert.Equal(t, andLset(a, b), []uint32(buf.Slice(nil)))
		assert.Equal(t, buf.Len(), len(buf.Slice(nil)))
	}
}

// randomSet returns a set whose containers range from empty and sparse to
// dense.
func randomSet(rnd *rand.Rand) *roaring.Set {
	s := roaring.New()
	for key := range uint32(6) {
		density := []uint32{0, 20, 2000, 6000, 50000}[rnd.IntN(5)]
		for range density {
			s.Add(key<<16 | rnd.Uint32N(1<<16))
		}
	}
	return s
}

// andLset computes a ∧ b with lset, sorted, as the expected result.
func andLset(a, b *roaring.Set) []uint32 {
	inB := lset.New(b.Slice(nil)...)
	var out []uint32
	for v := range a.Iter().Seq {
		if inB.Contains(v) {
			out = append(out, v)
		}
	}
	return out
}
