package liter

type Nested[Key, Val any] struct {
	Keys   Iter[Key]
	Vals   Iter[Val]
	Lookup func(Key) Iter[Val]
	I      int
}

func NewNested[Key, Val any](keys Iter[Key], lookup func(Key) Iter[Val]) *Nested[Key, Val] {
	return &Nested[Key, Val]{
		Keys:   keys,
		Lookup: lookup,
	}
}

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

func (n *Nested[Key, Val]) Done() bool {
	return n.Keys.Done()
}

func (n *Nested[Key, Val]) Idx() int {
	return n.I
}

func (n *Nested[Key, Val]) Wrap() Wrapper[Val] {
	return Wrapper[Val]{n}
}
