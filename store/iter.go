package store

import "github.com/adamcolton/luce/util/filter"

type Iter struct {
	FlatStore
	Key []byte
	filter.Filter[[]byte]
	I int
}

func NewIter(s FlatStore) *Iter {
	return &Iter{
		FlatStore: s,
		I:         -1,
	}
}

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

func (i *Iter) Cur() (key []byte, done bool) {
	if i.I == -1 {
		return i.Next()
	}
	return i.Key, i.Done()
}

func (i *Iter) Done() bool {
	return i.Key == nil && i.I > -1
}

func (i *Iter) Idx() int {
	return i.I
}

func (i *Iter) CurVal() (key []byte, r Record, done bool) {
	key, done = i.Cur()
	if done {
		return
	}
	r = i.FlatStore.Get(key)
	return
}
