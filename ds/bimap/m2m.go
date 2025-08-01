package bimap

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/flow"
)

// M2M holds a Many-to-Many mapping
type M2M[A, B comparable] struct {
	a2b map[A]*lset.Set[B]
	b2a map[B]*lset.Set[A]
}

// NewM2M creates a many-to-many mapping
func NewM2M[A, B comparable]() *M2M[A, B] {
	return &M2M[A, B]{
		a2b: make(map[A]*lset.Set[B]),
		b2a: make(map[B]*lset.Set[A]),
	}
}

func newSet[T comparable]() *lset.Set[T] {
	return lset.New[T]()
}

// Add pair (a,b)
func (m2m *M2M[A, B]) Add(a A, b B) {
	m2m.a2b[a] = flow.NilCheck(m2m.a2b[a], newSet).Add(b)
	m2m.b2a[b] = flow.NilCheck(m2m.b2a[b], newSet).Add(a)
}

// Remove pair (a,b)
func (m2m *M2M[A, B]) Remove(a A, b B) {
	if sa := m2m.a2b[a]; sa != nil {
		sb := m2m.b2a[b]

		sa.Remove(b)
		if sa.Len() == 0 {
			delete(m2m.a2b, a)
		}

		sb.Remove(a)
		if sb.Len() == 0 {
			delete(m2m.b2a, b)
		}
	}
}

// Each visits every (a,b) pair in the many-to-many set
func (m2m *M2M[A, B]) Each(fn func(a A, b B, done *bool)) {
	var done bool
	for a, b := range m2m.a2b {
		b.Each(func(b B, innerDone *bool) {
			fn(a, b, &done)
			*innerDone = done
		})
		if done {
			return
		}
	}
}

// A returns the set of all B values that have a mapping to the given value.
func (m2m *M2M[A, B]) A(a A) lset.Reader[B] {
	s := m2m.a2b[a]
	if s == nil {
		return nil
	}
	return s
}

// B returns the set of all A values that have a mapping to the given value.
func (m2m *M2M[A, B]) B(b B) lset.Reader[A] {
	s := m2m.b2a[b]
	if s == nil {
		return nil
	}
	return s
}

// RemoveA removes the value and return the set of all B values that have a
// had a mapping to the given value.
func (m2m *M2M[A, B]) RemoveA(a A) *lset.Set[B] {
	out := m2m.a2b[a]
	delete(m2m.a2b, a)
	out.Each(func(b B, done *bool) {
		s := m2m.b2a[b].Remove(a)
		if s.Len() == 0 {
			delete(m2m.b2a, b)
		}
	})
	return out
}

// RemoveB removes the value and return the set of all A values that have a
// had a mapping to the given value.
func (m2m *M2M[A, B]) RemoveB(b B) *lset.Set[A] {
	out := m2m.b2a[b]
	delete(m2m.b2a, b)
	out.Each(func(a A, done *bool) {
		s := m2m.a2b[a].Remove(b)
		if s.Len() == 0 {
			delete(m2m.a2b, a)
		}
	})
	return out
}

// LenA returns the number of unique A values in the many-to-many map.
func (m2m *M2M[A, B]) LenA() int {
	return len(m2m.a2b)
}

// LenB returns the number of unique B values in the many-to-many map.
func (m2m *M2M[A, B]) LenB() int {
	return len(m2m.b2a)
}

// EachA calls fn for every unique A value in the many-to-many map.
func (m2m *M2M[A, B]) EachA(fn func(a A, b lset.Reader[B], done *bool)) {
	lmap.New(m2m.a2b).Each(func(a A, b *lset.Set[B], done *bool) {
		fn(a, b, done)
	})
}

// EachB calls fn for every unique B value in the many-to-many map.
func (m2m *M2M[A, B]) EachB(fn func(b B, a lset.Reader[A], done *bool)) {
	lmap.New(m2m.b2a).Each(func(b B, a *lset.Set[A], done *bool) {
		fn(b, a, done)
	})
}

// SliceA returns a Slice with every unique A value in the many-to-many map.
func (m2m *M2M[A, B]) SliceA(buf []A) slice.Slice[A] {
	return lmap.New(m2m.a2b).Keys(buf)
}

// SliceB returns a Slice with every unique B value in the many-to-many map.
func (m2m *M2M[A, B]) SliceB(buf []B) slice.Slice[B] {
	return lmap.New(m2m.b2a).Keys(buf)
}
