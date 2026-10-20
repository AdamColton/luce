package roaring

import (
	"math/bits"

	"github.com/adamcolton/luce/ds/slice"
)

const bitmapWords = 1 << 16 / 64

// Bitmap is a set of uint16 values stored as one bit per possible value. It
// always occupies 8KB, so it suits dense sets. The zero value is an empty set.
type Bitmap struct {
	words [bitmapWords]uint64
}

// Contains reports whether v is in the set.
func (b *Bitmap) Contains(v uint16) bool {
	return b.words[v>>6]&(1<<(v&63)) != 0
}

// Add values to the set. The Bitmap is returned to allow chaining.
func (b *Bitmap) Add(vs ...uint16) *Bitmap {
	for _, v := range vs {
		b.Insert(v)
	}
	return b
}

// Checksert reports whether v was already in the set. If it was not, it is
// added.
func (b *Bitmap) Checksert(v uint16) bool {
	return !b.Insert(v)
}

// Remove v from the set. It is not an error if v is not in the set. The Bitmap
// is returned to allow chaining.
func (b *Bitmap) Remove(v uint16) *Bitmap {
	b.Delete(v)
	return b
}

// Insert puts v in the set and reports whether it was newly added.
func (b *Bitmap) Insert(v uint16) bool {
	w, bit := &b.words[v>>6], uint64(1)<<(v&63)
	if *w&bit != 0 {
		return false
	}
	*w |= bit
	return true
}

// Delete takes v out of the set and reports whether it was present.
func (b *Bitmap) Delete(v uint16) bool {
	w, bit := &b.words[v>>6], uint64(1)<<(v&63)
	if *w&bit == 0 {
		return false
	}
	*w &^= bit
	return true
}

// Len returns the number of values in the set.
func (b *Bitmap) Len() int {
	n := 0
	for _, w := range b.words {
		n += bits.OnesCount64(w)
	}
	return n
}

// Array returns the set as an Array, using buf's memory when it has the
// capacity.
func (b *Bitmap) Array(buf []uint16) Array {
	values := slice.NewBuffer(buf).Empty(b.Len())
	for i, w := range b.words {
		// Each pass takes the lowest set bit and then clears it.
		for ; w != 0; w &= w - 1 {
			values = append(values, uint16(i<<6+bits.TrailingZeros64(w)))
		}
	}
	return Array{values: values}
}
