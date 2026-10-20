package roaring_test

import (
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestSetOr(t *testing.T) {
	// Merge two blocklists into one.
	ours := roaring.New(ip(10, 0, 0, 1), ip(192, 168, 1, 7))
	shared := roaring.New(ip(192, 168, 1, 7), ip(203, 0, 113, 9), ip(8, 8, 4, 4))

	all := ours.Or(shared, nil)
	assert.Equal(t, []uint32{ip(8, 8, 4, 4), ip(10, 0, 0, 1), ip(192, 168, 1, 7), ip(203, 0, 113, 9)}, []uint32(all.Slice(nil)))
	assert.Equal(t, 2, ours.Len(), "inputs are unchanged")
	assert.Equal(t, 3, shared.Len())

	// The result doesn't share memory with the inputs.
	all.Remove(ip(10, 0, 0, 1))
	assert.True(t, ours.Contains(ip(10, 0, 0, 1)))

	// Next week's merge reuses this week's result as the buffer.
	next := ours.Or(roaring.New(ip(1, 1, 1, 1)), all)
	assert.Same(t, all, next)
	assert.Equal(t, []uint32{ip(1, 1, 1, 1), ip(10, 0, 0, 1), ip(192, 168, 1, 7)}, []uint32(next.Slice(nil)))
}

func TestSetOrPairings(t *testing.T) {
	// Container 0 is Array in both, 1 is Array∨Bitmap, 2 is Bitmap∨Array and
	// 3 is Bitmap∨Bitmap. In 4 two Arrays together hold too many values for
	// an Array, so the result becomes a Bitmap. 5 and 6 are each only in one
	// set, one of each kind.
	a := roaring.New()
	fill(a, 0, 0, 1000, 2)
	fill(a, 1, 0, 1000, 3)
	fill(a, 2, 0, 60000, 7)
	fill(a, 3, 0, 60000, 2)
	fill(a, 4, 0, 3000, 1)
	fill(a, 5, 0, 60000, 1)

	b := roaring.New()
	fill(b, 0, 0, 1000, 3)
	fill(b, 1, 0, 60000, 2)
	fill(b, 2, 0, 1000, 5)
	fill(b, 3, 0, 60000, 3)
	fill(b, 4, 2000, 5000, 1)
	fill(b, 6, 0, 10, 1)

	got := a.Or(b, nil)
	assert.Equal(t, orLset(a, b), []uint32(got.Slice(nil)))
	assert.Equal(t, got.Len(), len(got.Slice(nil)))
	assert.True(t, got.Contains(4<<16|4999))
	assert.Equal(t, got.Slice(nil), b.Or(a, nil).Slice(nil), "Or is symmetric")

	// Copied containers are independent of the inputs, for both kinds.
	got.Remove(5 << 16).Remove(6 << 16)
	assert.True(t, a.Contains(5<<16))
	assert.True(t, b.Contains(6<<16))
}

func TestSetOrEmpty(t *testing.T) {
	var empty roaring.Set
	a := roaring.New(1, 2, 3)
	assert.Equal(t, []uint32{1, 2, 3}, []uint32(a.Or(&empty, nil).Slice(nil)))
	assert.Equal(t, []uint32{1, 2, 3}, []uint32(empty.Or(a, nil).Slice(nil)))
	assert.Equal(t, 0, empty.Or(&empty, nil).Len())
}

func TestSetOrRandom(t *testing.T) {
	rnd := rand.New(rand.NewPCG(5, 6))
	var buf *roaring.Set
	for range 20 {
		a, b := randomSet(rnd), randomSet(rnd)
		buf = a.Or(b, buf)
		assert.Equal(t, orLset(a, b), []uint32(buf.Slice(nil)))
		assert.Equal(t, buf.Len(), len(buf.Slice(nil)))
	}
}

// orLset computes a ∨ b with lset, sorted, as the expected result.
func orLset(a, b *roaring.Set) []uint32 {
	out := lset.New(a.Slice(nil)...).Add(b.Slice(nil)...).Slice(nil)
	slices.Sort(out)
	return out
}
