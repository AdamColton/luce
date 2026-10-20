package roaring

import "slices"

// Array is a set of uint16 values kept as a sorted slice. It uses 2 bytes per
// value, so it suits sparse sets. The zero value is an empty set.
type Array struct {
	values []uint16
}

// Contains reports whether v is in the set.
func (a *Array) Contains(v uint16) bool {
	_, found := slices.BinarySearch(a.values, v)
	return found
}

// Add values to the set. The Array is returned to allow chaining.
func (a *Array) Add(vs ...uint16) *Array {
	for _, v := range vs {
		a.Insert(v)
	}
	return a
}

// Checksert reports whether v was already in the set. If it was not, it is
// added.
func (a *Array) Checksert(v uint16) bool {
	return !a.Insert(v)
}

// Remove v from the set. It is not an error if v is not in the set. The Array
// is returned to allow chaining.
func (a *Array) Remove(v uint16) *Array {
	a.Delete(v)
	return a
}

// Insert puts v in the set and reports whether it was newly added.
func (a *Array) Insert(v uint16) bool {
	i, found := slices.BinarySearch(a.values, v)
	if found {
		return false
	}
	a.values = slices.Insert(a.values, i, v)
	return true
}

// Delete takes v out of the set and reports whether it was present.
func (a *Array) Delete(v uint16) bool {
	i, found := slices.BinarySearch(a.values, v)
	if !found {
		return false
	}
	a.values = slices.Delete(a.values, i, i+1)
	return true
}

// Len returns the number of values in the set.
func (a *Array) Len() int {
	return len(a.values)
}

// Bitmap returns the set as a Bitmap. If buf is not nil it is cleared and
// used for the result.
func (a *Array) Bitmap(buf *Bitmap) *Bitmap {
	if buf == nil {
		buf = &Bitmap{}
	} else {
		*buf = Bitmap{}
	}
	for _, v := range a.values {
		buf.Insert(v)
	}
	return buf
}
