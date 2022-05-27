package huffman

// Lookup is used during encoding to finds the bit representation for a value in
// the tree. If T is comparable use NewLookup to create a Lookup. If T is not
// comparable, then use NewTranslateLookup and provide a function to translate
// from T to a comparable type.
type Lookup[T any] interface {
	Get(v T) *Bits
	All() []T
}

// Encode data to bits using the lookup. Calling Tree.ReadAll on these bits will
// recover the original data.
func Encode[T any](data []T, l Lookup[T]) *Bits {
	b := &Bits{}
	for _, d := range data {
		b.WriteBits(l.Get(d))
	}
	return b.Reset()
}

type mapLookup[T comparable] map[T]*Bits

// NewLookup creates a lookup on a Tree with a comparable type.
func NewLookup[T comparable](t Tree[T]) Lookup[T] {
	n := t.(*huffNode[T])
	l := make(mapLookup[T])
	l.insert(n, &Bits{})
	return l
}

func (l mapLookup[T]) insert(n *huffNode[T], b *Bits) {
	if n.branch[0] == nil {
		l[n.v] = b.Reset()
		return
	}
	l.insert(n.branch[0], b.Copy().Write(0))
	l.insert(n.branch[1], b.Write(1))
}

func (l mapLookup[T]) Get(v T) *Bits {
	return l[v]
}

func (l mapLookup[T]) All() []T {
	out := make([]T, 0, len(l))
	for t := range l {
		out = append(out, t)
	}
	return out
}

type translateLookup[K comparable, T any] struct {
	table    map[K]*Bits
	all      []T
	keyMaker func(T) K
}

func (l *translateLookup[K, T]) Get(v T) *Bits {
	return l.table[l.keyMaker(v)]
}

func (l *translateLookup[K, T]) All() []T {
	return l.all
}

// NewTranslateLookup creates a lookup when T is not comparable. A translator
// function must be provided.
func NewTranslateLookup[K comparable, T any](t Tree[T], translator func(T) K) Lookup[T] {
	n := t.(*huffNode[T])
	l := &translateLookup[K, T]{
		table:    make(map[K]*Bits),
		keyMaker: translator,
	}
	l.insert(n, &Bits{})
	return l
}

func (l *translateLookup[K, T]) insert(n *huffNode[T], b *Bits) {
	if n.branch[0] == nil {
		l.all = append(l.all, n.v)
		l.table[l.keyMaker(n.v)] = b.Reset()
		return
	}
	l.insert(n.branch[0], b.Copy().Write(0))
	l.insert(n.branch[1], b.Write(1))
}
