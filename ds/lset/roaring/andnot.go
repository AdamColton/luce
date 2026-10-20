package roaring

import "math/bits"

// AndNot returns the values that are in s but not in other. If buf is not nil,
// it is cleared and used for the result, reusing its memory; it must not be s
// or other.
func (s *Set) AndNot(other *Set, buf *Set) *Set {
	out, p := reuse(buf)
	// Walk s's entries, advancing through other's, which are also sorted by
	// key. An entry with no match in other is copied whole.
	j := 0
	for i := range s.entries {
		a := &s.entries[i]
		for j < len(other.entries) && other.entries[j].key < a.key {
			j++
		}
		e := entry{key: a.key}
		if j < len(other.entries) && other.entries[j].key == a.key {
			e.c, e.n = andNot(a, &other.entries[j], p)
		} else {
			e.c, e.n = clone(a.c, p), a.n
		}
		if e.n == 0 {
			p.put(e.c)
			continue
		}
		out.entries = append(out.entries, e)
		out.n += e.n
	}
	return out
}

// andNot returns the values in a's container that are not in b's, and their
// count.
func andNot(a, b *entry, p *pool) (container, int) {
	ba, aIsBitmap := a.c.(*Bitmap)
	bb, bIsBitmap := b.c.(*Bitmap)
	switch {
	case aIsBitmap && bIsBitmap:
		return andNotBitmaps(ba, bb, p)
	case aIsBitmap:
		return bitmapAndNotArray(ba, b.c.(*Array), a.n, p)
	case bIsBitmap:
		return arrayAndNotBitmap(a.c.(*Array), bb, p.array())
	}
	return andNotArrays(a.c.(*Array), b.c.(*Array), p.array())
}

// andNotArrays merges two sorted slices, keeping the values only in a.
func andNotArrays(a, b, out *Array) (container, int) {
	j := 0
	for _, v := range a.values {
		for j < len(b.values) && b.values[j] < v {
			j++
		}
		if j == len(b.values) || b.values[j] != v {
			out.values = append(out.values, v)
		}
	}
	return out, len(out.values)
}

func arrayAndNotBitmap(a *Array, b *Bitmap, out *Array) (container, int) {
	for _, v := range a.values {
		if !b.Contains(v) {
			out.values = append(out.values, v)
		}
	}
	return out, len(out.values)
}

// bitmapAndNotArray copies the Bitmap, which holds n values, and removes the
// Array's values from it.
func bitmapAndNotArray(a *Bitmap, b *Array, n int, p *pool) (container, int) {
	out := p.bitmap()
	*out = *a
	for _, v := range b.values {
		if out.Delete(v) {
			n--
		}
	}
	return fit(out, n, p), n
}

func andNotBitmaps(a, b *Bitmap, p *pool) (container, int) {
	out := p.bitmap()
	n := 0
	for i := range out.words {
		w := a.words[i] &^ b.words[i]
		out.words[i] = w
		n += bits.OnesCount64(w)
	}
	return fit(out, n, p), n
}

// fit returns b, which holds n values, or b converted to an Array if n is few
// enough for one.
func fit(b *Bitmap, n int, p *pool) container {
	if n > arrayMax {
		return b
	}
	arr := p.array()
	*arr = b.Array(arr.values)
	p.put(b)
	return arr
}
