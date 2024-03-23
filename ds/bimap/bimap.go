package bimap

type Bimap[A, B comparable] struct {
	a2b map[A]B
	b2a map[B]A
}

func New[A, B comparable](size int) *Bimap[A, B] {
	return &Bimap[A, B]{
		a2b: make(map[A]B, size),
		b2a: make(map[B]A, size),
	}
}

type Deleted[X comparable] struct {
	Value   X
	Deleted bool
}

type Bidelete[A, B comparable] struct {
	A Deleted[A]
	B Deleted[B]
}

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

func (bi *Bimap[A, B]) DeleteA(a A) (d Deleted[B]) {
	d.Value, d.Deleted = bi.a2b[a]
	if d.Deleted {
		delete(bi.a2b, a)
		delete(bi.b2a, d.Value)
	}

	return
}

func (bi *Bimap[A, B]) DeleteB(b B) (d Deleted[A]) {
	d.Value, d.Deleted = bi.b2a[b]
	if d.Deleted {
		delete(bi.b2a, b)
		delete(bi.a2b, d.Value)
	}

	return
}

func (bi *Bimap[A, B]) A(a A) (b B, found bool) {
	b, found = bi.a2b[a]
	return
}

func (bi *Bimap[A, B]) B(b B) (a A, found bool) {
	a, found = bi.b2a[b]
	return
}
