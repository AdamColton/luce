package store

import "github.com/adamcolton/luce/util/filter"

type Iter struct {
	Store
	Key []byte
	filter.Filter[[]byte]
	I int
}

func NewIter(s Store) *Iter {
	return &Iter{
		Store: s,
		I:     -1,
	}
}

func (i *Iter) Next() (key []byte, done bool) {
	i.I++
	for {
		i.Key = i.Store.Next(i.Key)
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
	r = i.Store.Get(key)
	return
}
