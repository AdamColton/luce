package liter

// Nested takes an Iter of keys and a function to use a key to get an Iter of
// values. It iterates over the values calling the Lookup each time the Vals
// iterator is done until the Keys are done. Nested itself fulfills Iter[Val].
type Nested[Key, Val any] struct {
	Keys   Iter[Key]
	Vals   Iter[Val]
	Lookup func(Key) Iter[Val]
	I      int
}

// NewNested creates an instance of Nested from the keys and lookup function.
func NewNested[Key, Val any](keys Iter[Key], lookup func(Key) Iter[Val]) *Nested[Key, Val] {
	return &Nested[Key, Val]{
		Keys:   keys,
		Lookup: lookup,
	}
}

// Next fulfills Iter. It returns the next value and a bool indicating if the
// iterator is done.
func (n *Nested[Key, Val]) Next() (val Val, done bool) {
	if n.Vals == nil || n.Vals.Done() {
		var k Key
		k, done = n.Keys.Next()
		if done {
			return
		}
		n.Vals = n.Lookup(k)
		val, done = n.Vals.Cur()
	} else {
		val, done = n.Vals.Next()
	}
	if done {
		return n.Next()
	}
	n.I++
	return
}

// Cur fulfills Iter. It returns the current value and a bool indicating if the
// iterator is done.
func (n *Nested[Key, Val]) Cur() (val Val, done bool) {
	if n.Vals == nil {
		k, keysDone := n.Keys.Cur()
		for {
			if keysDone {
				done = true
				return
			}
			n.Vals = n.Lookup(k)
			val, done = n.Vals.Cur()
			if !done {
				return
			}
			k, keysDone = n.Keys.Next()
		}
	}
	val, _ = n.Vals.Cur()
	return val, n.Keys.Done()
}

// Done fulfills Iter and returns true if the iterator is done.
func (n *Nested[Key, Val]) Done() bool {
	return n.Keys.Done()
}

// Idx fulfills Iter and returns the index of the current value.
func (n *Nested[Key, Val]) Idx() int {
	return n.I
}

// Wrap the instance of Nested for all the features the Wrapper provides.
func (n *Nested[Key, Val]) Wrap() Wrapper[Val] {
	return Wrapper[Val]{n}
}
