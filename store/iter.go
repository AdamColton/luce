package store

import "github.com/adamcolton/luce/util/filter"

// Iter walks the keys of a FlatStore in order. It fulfills liter.Iter. Set
// Filter to skip keys it returns false for. Iter embeds the FlatStore, but its
// own Next hides the store's Next.
type Iter struct {
	FlatStore
	// Key is the current key, nil before the first Next and after the last.
	Key []byte
	filter.Filter[[]byte]
	// I is the index of the current key, -1 before the first Next.
	I int
}

// NewIter creates an Iter that has not started. The first Next or Cur moves to
// the first key.
func NewIter(s FlatStore) *Iter {
	return &Iter{
		FlatStore: s,
		I:         -1,
	}
}

// Next moves to the next key that passes the Filter and returns it. done is
// true once there are no more keys.
func (i *Iter) Next() (key []byte, done bool) {
	i.I++
	for {
		i.Key = i.FlatStore.Next(i.Key)
		if i.Key == nil || i.Filter == nil || i.Filter(i.Key) {
			break
		}
	}
	return i.Key, i.Done()
}

// Cur returns the current key. If the Iter has not started it moves to the
// first key.
func (i *Iter) Cur() (key []byte, done bool) {
	if i.I == -1 {
		return i.Next()
	}
	return i.Key, i.Done()
}

// Done is true once the Iter has moved past the last key.
func (i *Iter) Done() bool {
	return i.Key == nil && i.I > -1
}

// Idx is the index of the current key, counting only the keys that passed the
// Filter. It is -1 before the first Next.
func (i *Iter) Idx() int {
	return i.I
}

// CurVal returns the current key and the Record stored at it. The Record is
// empty when done.
func (i *Iter) CurVal() (key []byte, r Record, done bool) {
	key, done = i.Cur()
	if done {
		return
	}
	r = i.FlatStore.Get(key)
	return
}
