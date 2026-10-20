package roaring

import "math/bits"

// Or returns the values that are in s, other or both. If buf is not nil, it is
// cleared and used for the result, reusing its memory; it must not be s or
// other.
func (s *Set) Or(other *Set, buf *Set) *Set {
	out, p := reuse(buf)
	// Merge the two entry lists, which are sorted by key. An entry in only
	// one set is copied, so the result never shares memory with an input.
	i, j := 0, 0
	for i < len(s.entries) || j < len(other.entries) {
		var e entry
		switch {
		case j == len(other.entries) || (i < len(s.entries) && s.entries[i].key < other.entries[j].key):
			e = s.entries[i]
			e.c = clone(e.c, p)
			i++
		case i == len(s.entries) || s.entries[i].key > other.entries[j].key:
			e = other.entries[j]
			e.c = clone(e.c, p)
			j++
		default:
			e.key = s.entries[i].key
			e.c, e.n = or(&s.entries[i], &other.entries[j], p)
			i++
			j++
		}
		out.entries = append(out.entries, e)
		out.n += e.n
	}
	return out
}

// clone copies c into a container from the pool.
func clone(c container, p *pool) container {
	if b, ok := c.(*Bitmap); ok {
		out := p.bitmap()
		*out = *b
		return out
	}
	out := p.array()
	out.values = append(out.values, c.(*Array).values...)
	return out
}

// or returns the union of two entries' containers and its count.
func or(a, b *entry, p *pool) (container, int) {
	ba, aIsBitmap := a.c.(*Bitmap)
	bb, bIsBitmap := b.c.(*Bitmap)
	switch {
	case aIsBitmap && bIsBitmap:
		return orBitmaps(ba, bb, p.bitmap())
	case aIsBitmap:
		return orArrayBitmap(b.c.(*Array), ba, a.n, p.bitmap())
	case bIsBitmap:
		return orArrayBitmap(a.c.(*Array), bb, b.n, p.bitmap())
	}
	return orArrays(a.c.(*Array), b.c.(*Array), p)
}

// orArrays merges two sorted slices, keeping the values in either. The result
// can hold too many values for an Array, in which case it is converted.
func orArrays(a, b *Array, p *pool) (container, int) {
	out := p.array()
	i, j := 0, 0
	for i < len(a.values) && j < len(b.values) {
		switch av, bv := a.values[i], b.values[j]; {
		case av < bv:
			out.values = append(out.values, av)
			i++
		case av > bv:
			out.values = append(out.values, bv)
			j++
		default:
			out.values = append(out.values, av)
			i++
			j++
		}
	}
	out.values = append(out.values, a.values[i:]...)
	out.values = append(out.values, b.values[j:]...)
	n := len(out.values)
	if n <= arrayMax {
		return out, n
	}
	bm := out.Bitmap(p.bitmap())
	p.put(out)
	return bm, n
}

// orArrayBitmap copies the Bitmap, which holds n values, and adds the Array's
// values to it.
func orArrayBitmap(a *Array, b *Bitmap, n int, out *Bitmap) (container, int) {
	*out = *b
	for _, v := range a.values {
		if out.Insert(v) {
			n++
		}
	}
	return out, n
}

func orBitmaps(a, b, out *Bitmap) (container, int) {
	n := 0
	for i := range out.words {
		w := a.words[i] | b.words[i]
		out.words[i] = w
		n += bits.OnesCount64(w)
	}
	return out, n
}
