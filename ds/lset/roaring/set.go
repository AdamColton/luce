package roaring

import (
	"cmp"
	"slices"
)

// arrayMax is the most values a container holds as an Array; above it the
// container is a Bitmap. At 4096 values both take 8KB.
const arrayMax = 4096

// Set is a set of uint32 values stored as a Roaring bitmap. The zero value is
// an empty set. A nil *Set can be used with Len (0); any other method panics
// on a nil Set.
type Set struct {
	entries []entry // sorted by key
	n       int
}

// entry holds the values whose high 16 bits are key, as their low 16 bits.
type entry struct {
	key uint16
	n   int // count of c, so changes don't need Bitmap.Len
	c   container
}

// container is what Array and Bitmap have in common.
type container interface {
	Contains(v uint16) bool
	Insert(v uint16) bool
	Delete(v uint16) bool
}

// New creates a set containing the provided values. Repeated values are only
// in the set once.
func New(vs ...uint32) *Set {
	return (&Set{}).Add(vs...)
}

// Contains reports whether v is in the set.
func (s *Set) Contains(v uint32) bool {
	i, found := s.find(uint16(v >> 16))
	return found && s.entries[i].c.Contains(uint16(v))
}

// Add values to the set. The Set is returned to allow chaining.
func (s *Set) Add(vs ...uint32) *Set {
	for _, v := range vs {
		s.Insert(v)
	}
	return s
}

// Checksert reports whether v was already in the set. If it was not, it is
// added.
func (s *Set) Checksert(v uint32) bool {
	return !s.Insert(v)
}

// Remove v from the set. It is not an error if v is not in the set. The Set is
// returned to allow chaining.
func (s *Set) Remove(v uint32) *Set {
	s.Delete(v)
	return s
}

// Len returns the number of values in the set. If the Set is nil, it is 0.
func (s *Set) Len() int {
	if s == nil {
		return 0
	}
	return s.n
}

// Insert puts v in the set and reports whether it was newly added.
func (s *Set) Insert(v uint32) bool {
	i, found := s.find(uint16(v >> 16))
	if !found {
		s.entries = slices.Insert(s.entries, i, entry{key: uint16(v >> 16), c: &Array{}})
	}
	e := &s.entries[i]
	if !e.c.Insert(uint16(v)) {
		return false
	}
	e.n++
	s.n++
	if e.n == arrayMax+1 {
		e.c = e.c.(*Array).Bitmap(nil)
	}
	return true
}

// Delete takes v out of the set and reports whether it was present.
func (s *Set) Delete(v uint32) bool {
	i, found := s.find(uint16(v >> 16))
	if !found {
		return false
	}
	e := &s.entries[i]
	if !e.c.Delete(uint16(v)) {
		return false
	}
	e.n--
	s.n--
	switch e.n {
	case 0:
		s.entries = slices.Delete(s.entries, i, i+1)
	case arrayMax:
		a := e.c.(*Bitmap).Array(nil)
		e.c = &a
	}
	return true
}

func (s *Set) find(key uint16) (int, bool) {
	return slices.BinarySearchFunc(s.entries, key, func(e entry, key uint16) int {
		return cmp.Compare(e.key, key)
	})
}
