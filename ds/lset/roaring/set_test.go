package roaring_test

import (
	"math/rand/v2"
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

// ip packs an IPv4 address into a uint32.
func ip(a, b, c, d byte) uint32 {
	return uint32(a)<<24 | uint32(b)<<16 | uint32(c)<<8 | uint32(d)
}

func TestSet(t *testing.T) {
	// A blocklist of IPv4 addresses. Each /16 network gets its own container.
	blocked := roaring.New(ip(10, 0, 0, 1), ip(192, 168, 1, 7), ip(203, 0, 113, 9))
	assert.Equal(t, 3, blocked.Len())
	assert.True(t, blocked.Contains(ip(192, 168, 1, 7)))
	assert.False(t, blocked.Contains(ip(192, 168, 1, 8)), "same /16, not blocked")
	assert.False(t, blocked.Contains(ip(8, 8, 8, 8)), "no container for this /16")

	assert.True(t, blocked.Checksert(ip(10, 0, 0, 1)), "already blocked")
	assert.False(t, blocked.Checksert(ip(10, 0, 0, 2)), "newly blocked")
	assert.Equal(t, 4, blocked.Len())

	blocked.Remove(ip(203, 0, 113, 9)).Remove(ip(8, 8, 8, 8))
	assert.False(t, blocked.Contains(ip(203, 0, 113, 9)))
	assert.Equal(t, 3, blocked.Len())

	// Insert and Delete report whether anything changed.
	assert.True(t, blocked.Insert(ip(198, 51, 100, 4)))
	assert.False(t, blocked.Insert(ip(198, 51, 100, 4)), "already blocked")
	assert.True(t, blocked.Delete(ip(198, 51, 100, 4)))
	assert.False(t, blocked.Delete(ip(198, 51, 100, 4)), "already unblocked")

	// Block all of 172.16.0.0/16 except 172.16.255.255. That is far more than
	// 4096 addresses, so the network's container becomes a Bitmap.
	for c := range 256 {
		for d := range 256 {
			blocked.Add(ip(172, 16, byte(c), byte(d)))
		}
	}
	blocked.Remove(ip(172, 16, 255, 255))
	assert.Equal(t, 3+65535, blocked.Len())
	assert.True(t, blocked.Contains(ip(172, 16, 40, 2)))
	assert.False(t, blocked.Contains(ip(172, 16, 255, 255)))
	assert.False(t, blocked.Contains(ip(172, 17, 0, 0)))

	// Unblock most of it again; the container shrinks back to an Array.
	for c := range 255 {
		for d := range 256 {
			blocked.Remove(ip(172, 16, byte(c), byte(d)))
		}
	}
	assert.Equal(t, 3+255, blocked.Len())
	assert.True(t, blocked.Contains(ip(172, 16, 255, 0)))
	assert.False(t, blocked.Contains(ip(172, 16, 40, 2)))
}

func TestSetChecksertOnce(t *testing.T) {
	// Checksert as a "do this once" test: log each address the first time it
	// connects, however many times it connects.
	connections := []uint32{
		ip(10, 0, 0, 1), ip(10, 0, 0, 2), ip(10, 0, 0, 1), ip(10, 0, 0, 1), ip(10, 0, 0, 3),
	}
	var seen roaring.Set
	var firstSeen []uint32
	for _, addr := range connections {
		if !seen.Checksert(addr) {
			firstSeen = append(firstSeen, addr)
		}
	}
	assert.Equal(t, []uint32{ip(10, 0, 0, 1), ip(10, 0, 0, 2), ip(10, 0, 0, 3)}, firstSeen)
}

func TestSetNil(t *testing.T) {
	var s *roaring.Set
	assert.Equal(t, 0, s.Len())
}

func TestSetZeroValue(t *testing.T) {
	var s roaring.Set
	assert.False(t, s.Contains(0))
	s.Remove(0)
	s.Add(0, 1<<32-1)
	assert.True(t, s.Contains(0))
	assert.True(t, s.Contains(1<<32-1))
	assert.Equal(t, 2, s.Len())
}

func TestSetMatchesLset(t *testing.T) {
	// Random changes concentrated in a few containers, so they cross the
	// Array/Bitmap threshold in both directions, checked against lset.Set.
	rnd := rand.New(rand.NewPCG(1, 2))
	r := roaring.New()
	l := lset.New[uint32]()
	for range 200_000 {
		v := rnd.Uint32N(3)<<16 | rnd.Uint32N(12000)
		if rnd.IntN(5) < 3 {
			assert.Equal(t, l.Contains(v), r.Checksert(v))
			l.Add(v)
		} else {
			r.Remove(v)
			l.Remove(v)
		}
	}
	assert.Equal(t, l.Len(), r.Len())
	for hi := range uint32(4) {
		for lo := range uint32(12000) {
			v := hi<<16 | lo
			if l.Contains(v) != r.Contains(v) {
				t.Fatalf("disagree on %d", v)
			}
		}
	}
}
