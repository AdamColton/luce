package bimap

// Bimap allows for bidirectional lookup: each A is paired with exactly one B and
// each B with exactly one A. It is not safe for concurrent use.
type Bimap[A, B comparable] struct {
	a2b map[A]B
	b2a map[B]A
}

// New creates a bimap with room for size pairs.
func New[A, B comparable](size int) *Bimap[A, B] {
	return &Bimap[A, B]{
		a2b: make(map[A]B, size),
		b2a: make(map[B]A, size),
	}
}

// Deleted is returned when a key is deleted from the Bimap. Value is the partner
// of the deleted key and Deleted is true if the key was in the Bimap.
type Deleted[X comparable] struct {
	Value   X
	Deleted bool
}

// Bidelete is returned when calling Add. It indicates if the insert caused
// values to be removed. A is the A value that was paired with the added b, and B
// is the B value that was paired with the added a.
type Bidelete[A, B comparable] struct {
	A Deleted[A]
	B Deleted[B]
}

// Add a pair of values to the Bimap. Any existing pair that uses a or b is
// removed first, and the returned Bidelete indicates which values were removed.
// Adding a pair that is already in the Bimap reports the B as deleted, although
// the Bimap is unchanged.
func (bi *Bimap[A, B]) Add(a A, b B) (bd Bidelete[A, B]) {
	if d := bi.DeleteA(a); d.Deleted {
		bd.B = d
	}
	if d := bi.DeleteB(b); d.Deleted {
		bd.A = d
	}

	bi.a2b[a] = b
	bi.b2a[b] = a

	return
}

// DeleteA deletes a keypair by its A value. The returned Deleted holds the B
// that was paired with it.
func (bi *Bimap[A, B]) DeleteA(a A) (d Deleted[B]) {
	d.Value, d.Deleted = bi.a2b[a]
	if d.Deleted {
		delete(bi.a2b, a)
		delete(bi.b2a, d.Value)
	}

	return
}

// DeleteB deletes a keypair by its B value. The returned Deleted holds the A
// that was paired with it.
func (bi *Bimap[A, B]) DeleteB(b B) (d Deleted[A]) {
	d.Value, d.Deleted = bi.b2a[b]
	if d.Deleted {
		delete(bi.b2a, b)
		delete(bi.a2b, d.Value)
	}

	return
}

// A does a lookup using an 'A' key and returns its B.
func (bi *Bimap[A, B]) A(a A) (b B, found bool) {
	b, found = bi.a2b[a]
	return
}

// B does a lookup using a 'B' key and returns its A.
func (bi *Bimap[A, B]) B(b B) (a A, found bool) {
	a, found = bi.b2a[b]
	return
}

// Each invokes the provided function for each keypair, in no particular order.
// Setting done to true stops the iteration.
func (bi *Bimap[A, B]) Each(fn func(a A, b B, done *bool)) {
	done := false
	for a, b := range bi.a2b {
		fn(a, b, &done)
		if done {
			break
		}
	}
}
