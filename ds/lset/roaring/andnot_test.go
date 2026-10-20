package roaring_test

import (
	"math/rand/v2"
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestSetAndNot(t *testing.T) {
	// Which of today's connections came from addresses that aren't blocked?
	today := roaring.New(ip(192, 168, 1, 7), ip(8, 8, 8, 8), ip(203, 0, 113, 9), ip(10, 0, 0, 2))
	blocked := roaring.New(ip(10, 0, 0, 1), ip(192, 168, 1, 7), ip(203, 0, 113, 9))

	allowed := today.AndNot(blocked, nil)
	assert.Equal(t, []uint32{ip(8, 8, 8, 8), ip(10, 0, 0, 2)}, []uint32(allowed.Slice(nil)))
	assert.Equal(t, 4, today.Len(), "inputs are unchanged")
	assert.Equal(t, 3, blocked.Len())

	// The result doesn't share memory with the inputs.
	allowed.Remove(ip(8, 8, 8, 8))
	assert.True(t, today.Contains(ip(8, 8, 8, 8)))

	// Tomorrow's check reuses today's result as the buffer.
	tomorrow := roaring.New(ip(10, 0, 0, 1), ip(1, 1, 1, 1))
	allowed2 := tomorrow.AndNot(blocked, allowed)
	assert.Same(t, allowed, allowed2)
	assert.Equal(t, []uint32{ip(1, 1, 1, 1)}, []uint32(allowed2.Slice(nil)))
}

func TestSetAndNotPairings(t *testing.T) {
	// Container 0 is Array−Array, 1 is Array−Bitmap, 2 is Bitmap−Array with
	// the result still a Bitmap, and 3 is Bitmap−Bitmap with the result still
	// a Bitmap. In 4 (Bitmap−Array) and 5 (Bitmap−Bitmap) most values are
	// removed, so the results become Arrays. In 6 everything is removed. 7 is
	// only in a, and 8 only in b.
	a := roaring.New()
	fill(a, 0, 0, 1000, 2)
	fill(a, 1, 0, 1000, 3)
	fill(a, 2, 0, 60000, 7)
	fill(a, 3, 0, 60000, 2)
	fill(a, 4, 0, 5000, 1)
	fill(a, 5, 0, 60000, 1)
	fill(a, 6, 0, 100, 1)
	fill(a, 7, 0, 60000, 1)

	b := roaring.New()
	fill(b, 0, 0, 1000, 3)
	fill(b, 1, 0, 60000, 2)
	fill(b, 2, 0, 1000, 5)
	fill(b, 3, 0, 60000, 3)
	fill(b, 4, 0, 4000, 1)
	fill(b, 5, 10, 60000, 1)
	fill(b, 6, 0, 200, 1)
	fill(b, 8, 0, 10, 1)

	got := a.AndNot(b, nil)
	assert.Equal(t, andNotLset(a, b), []uint32(got.Slice(nil)))
	assert.Equal(t, got.Len(), len(got.Slice(nil)))
	assert.True(t, got.Contains(4<<16|4500))
	assert.False(t, got.Contains(6<<16|5))
	assert.False(t, got.Contains(8<<16|5))

	// A copied container is independent of the input.
	got.Remove(7 << 16)
	assert.True(t, a.Contains(7<<16))
}

func TestSetAndNotEmpty(t *testing.T) {
	var empty roaring.Set
	a := roaring.New(1, 2, 3)
	assert.Equal(t, []uint32{1, 2, 3}, []uint32(a.AndNot(&empty, nil).Slice(nil)))
	assert.Equal(t, 0, empty.AndNot(a, nil).Len())
	assert.Equal(t, 0, a.AndNot(a, nil).Len())
}

func TestSetAndNotRandom(t *testing.T) {
	rnd := rand.New(rand.NewPCG(7, 8))
	var buf *roaring.Set
	for range 20 {
		a, b := randomSet(rnd), randomSet(rnd)
		buf = a.AndNot(b, buf)
		assert.Equal(t, andNotLset(a, b), []uint32(buf.Slice(nil)))
		assert.Equal(t, buf.Len(), len(buf.Slice(nil)))
	}
}

// andNotLset computes a − b with lset, sorted, as the expected result.
func andNotLset(a, b *roaring.Set) []uint32 {
	inB := lset.New(b.Slice(nil)...)
	var out []uint32
	for v := range a.Iter().Seq {
		if !inB.Contains(v) {
			out = append(out, v)
		}
	}
	return out
}
