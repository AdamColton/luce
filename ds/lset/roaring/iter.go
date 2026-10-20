package roaring

import (
	"math/bits"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
)

// Slice returns the values in the set in ascending order. The buffer is used
// if it has sufficient capacity. If the Set is nil, a nil slice is returned.
func (s *Set) Slice(buf []uint32) slice.Slice[uint32] {
	if s == nil {
		return nil
	}
	out := slice.NewBuffer(buf).Empty(s.n)
	for _, e := range s.entries {
		out = e.c.appendTo(out, uint32(e.key)<<16)
	}
	return out
}

// Iter returns an iterator over the values in the set in ascending order. The
// iterator also fulfills liter.Starter, and slice.Slicer, whose Slice returns
// the whole set regardless of the iterator's position. The set must not change
// while the iterator is in use. A nil Set gives an iterator that is done.
func (s *Set) Iter() liter.Wrapper[uint32] {
	if s == nil {
		s = &Set{}
	}
	i := &iter{s: s}
	i.Start()
	return liter.Wrap[uint32](i)
}

type iter struct {
	s   *Set
	ei  int // index of the current entry
	pos int // position in the current entry's container
	idx int
}

func (i *iter) Start() (uint32, bool) {
	i.ei, i.pos, i.idx = 0, 0, 0
	i.settle()
	return i.Cur()
}

func (i *iter) Next() (uint32, bool) {
	if !i.Done() {
		i.pos++
		i.idx++
		i.settle()
	}
	return i.Cur()
}

// settle moves to the first value at or after pos, moving on to later
// entries when the current one has no more.
func (i *iter) settle() {
	for ; i.ei < len(i.s.entries); i.ei, i.pos = i.ei+1, 0 {
		if pos, ok := i.s.entries[i.ei].c.nextFrom(i.pos); ok {
			i.pos = pos
			return
		}
	}
}

func (i *iter) Cur() (uint32, bool) {
	if i.Done() {
		return 0, true
	}
	e := &i.s.entries[i.ei]
	return uint32(e.key)<<16 | uint32(e.c.at(i.pos)), false
}

func (i *iter) Done() bool {
	return i.ei >= len(i.s.entries)
}

func (i *iter) Idx() int {
	return i.idx
}

func (i *iter) Slice(buf []uint32) []uint32 {
	return i.s.Slice(buf)
}

// A position in an Array is an index into values.

func (a *Array) nextFrom(i int) (int, bool) {
	return i, i < len(a.values)
}

func (a *Array) at(pos int) uint16 {
	return a.values[pos]
}

func (a *Array) appendTo(out []uint32, hi uint32) []uint32 {
	for _, v := range a.values {
		out = append(out, hi|uint32(v))
	}
	return out
}

// A position in a Bitmap is the value itself.

func (b *Bitmap) nextFrom(i int) (int, bool) {
	if i >= 1<<16 {
		return 0, false
	}
	wi := i >> 6
	w := b.words[wi] &^ (1<<(i&63) - 1) // drop the bits below i
	for {
		if w != 0 {
			return wi<<6 + bits.TrailingZeros64(w), true
		}
		wi++
		if wi == bitmapWords {
			return 0, false
		}
		w = b.words[wi]
	}
}

func (b *Bitmap) at(pos int) uint16 {
	return uint16(pos)
}

func (b *Bitmap) appendTo(out []uint32, hi uint32) []uint32 {
	for i, w := range b.words {
		// Each pass takes the lowest set bit and then clears it.
		for ; w != 0; w &= w - 1 {
			out = append(out, hi|uint32(i<<6+bits.TrailingZeros64(w)))
		}
	}
	return out
}
