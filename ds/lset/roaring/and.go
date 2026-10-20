package roaring

import "math/bits"

// And returns the values that are in both s and other. If buf is not nil, it
// is cleared and used for the result, reusing its memory; it must not be s or
// other.
func (s *Set) And(other *Set, buf *Set) *Set {
	out, p := reuse(buf)
	// Merge the two entry lists, which are sorted by key.
	i, j := 0, 0
	for i < len(s.entries) && j < len(other.entries) {
		a, b := &s.entries[i], &other.entries[j]
		switch {
		case a.key < b.key:
			i++
		case a.key > b.key:
			j++
		default:
			c, n := and(a.c, b.c, p)
			if n > 0 {
				out.entries = append(out.entries, entry{key: a.key, n: n, c: c})
				out.n += n
			} else {
				p.put(c)
			}
			i++
			j++
		}
	}
	return out
}

// reuse clears buf for use as a result, keeping its containers in a pool so
// their memory can be reused. A nil buf gives a new Set and an empty pool.
func reuse(buf *Set) (*Set, *pool) {
	p := &pool{}
	if buf == nil {
		return &Set{}, p
	}
	for _, e := range buf.entries {
		p.put(e.c)
	}
	buf.entries = buf.entries[:0]
	buf.n = 0
	return buf, p
}

// pool holds containers that are free to be reused.
type pool struct {
	arrays  []*Array
	bitmaps []*Bitmap
}

func (p *pool) put(c container) {
	switch c := c.(type) {
	case *Array:
		p.arrays = append(p.arrays, c)
	case *Bitmap:
		p.bitmaps = append(p.bitmaps, c)
	}
}

// array returns an empty Array.
func (p *pool) array() *Array {
	if len(p.arrays) == 0 {
		return &Array{}
	}
	a := p.arrays[len(p.arrays)-1]
	p.arrays = p.arrays[:len(p.arrays)-1]
	a.values = a.values[:0]
	return a
}

// bitmap returns a Bitmap whose contents are left over from earlier use; the
// caller must overwrite every word.
func (p *pool) bitmap() *Bitmap {
	if len(p.bitmaps) == 0 {
		return &Bitmap{}
	}
	b := p.bitmaps[len(p.bitmaps)-1]
	p.bitmaps = p.bitmaps[:len(p.bitmaps)-1]
	return b
}

// and returns the intersection of two containers and its count.
func and(a, b container, p *pool) (container, int) {
	ba, aIsBitmap := a.(*Bitmap)
	bb, bIsBitmap := b.(*Bitmap)
	switch {
	case aIsBitmap && bIsBitmap:
		return andBitmaps(ba, bb, p)
	case aIsBitmap:
		return andArrayBitmap(b.(*Array), ba, p.array())
	case bIsBitmap:
		return andArrayBitmap(a.(*Array), bb, p.array())
	}
	return andArrays(a.(*Array), b.(*Array), p.array())
}

// andArrays merges two sorted slices, keeping the values in both.
func andArrays(a, b, out *Array) (container, int) {
	i, j := 0, 0
	for i < len(a.values) && j < len(b.values) {
		switch av, bv := a.values[i], b.values[j]; {
		case av < bv:
			i++
		case av > bv:
			j++
		default:
			out.values = append(out.values, av)
			i++
			j++
		}
	}
	return out, len(out.values)
}

func andArrayBitmap(a *Array, b *Bitmap, out *Array) (container, int) {
	for _, v := range a.values {
		if b.Contains(v) {
			out.values = append(out.values, v)
		}
	}
	return out, len(out.values)
}

// andBitmaps ANDs word by word. The result can hold few enough values to
// belong in an Array, in which case it is converted.
func andBitmaps(a, b *Bitmap, p *pool) (container, int) {
	out := p.bitmap()
	n := 0
	for i := range out.words {
		w := a.words[i] & b.words[i]
		out.words[i] = w
		n += bits.OnesCount64(w)
	}
	if n > arrayMax {
		return out, n
	}
	arr := p.array()
	*arr = out.Array(arr.values)
	p.put(out)
	return arr, n
}
