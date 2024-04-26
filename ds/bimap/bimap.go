package bimap

// Bimap allows for bidirectional lookup.
type Bimap[A, B comparable] struct {
	a2b map[A]B
	b2a map[B]A
}

// New creates a bimap.
func New[A, B comparable](size int) *Bimap[A, B] {
	return &Bimap[A, B]{
		a2b: make(map[A]B, size),
		b2a: make(map[B]A, size),
	}
}

// Deleted is returned when a key from the Bimap. It is the partner to the
// deleted key.
type Deleted[X comparable] struct {
	Value   X
	Deleted bool
}

// Bidelete is returned when calling Add. It indicates if the insert caused
// values to be removed.
type Bidelete[A, B comparable] struct {
	A Deleted[A]
	B Deleted[B]
}

// Add a pair of values to the Bimap. The returned Bidelete indicates if any
// values were removed by the Add.
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

// DeleteA deletes a keypair by it's A value.
func (bi *Bimap[A, B]) DeleteA(a A) (d Deleted[B]) {
	d.Value, d.Deleted = bi.a2b[a]
	if d.Deleted {
		delete(bi.a2b, a)
		delete(bi.b2a, d.Value)
	}

	return
}

// DeleteB deletes a keypair by it's B value.
func (bi *Bimap[A, B]) DeleteB(b B) (d Deleted[A]) {
	d.Value, d.Deleted = bi.b2a[b]
	if d.Deleted {
		delete(bi.b2a, b)
		delete(bi.a2b, d.Value)
	}

	return
}

// A does a lookup using an 'A' key.
func (bi *Bimap[A, B]) A(a A) (b B, found bool) {
	b, found = bi.a2b[a]
	return
}

// B does a lookup using a 'B' key.
func (bi *Bimap[A, B]) B(b B) (a A, found bool) {
	a, found = bi.b2a[b]
	return
}

// Each invokes the provided function for each keypair.
func (bi *Bimap[A, B]) Each(fn func(a A, b B) bool) {
	for a, b := range bi.a2b {
		if fn(a, b) {
			break
		}
	}
}
