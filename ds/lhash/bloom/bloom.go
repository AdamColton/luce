package bloom

import (
	"github.com/adamcolton/luce/ds/lhash"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye"
)

const ErrOutOfRange = lerr.Str("Indexer returned index out of range")

type Filter[T any, U lhash.Uint] struct {
	bits *rye.Bits
	idx  lhash.Indexer[T, U]
}

func New[T any, U lhash.Uint](indexer lhash.Indexer[T, U], values ...T) *Filter[T, U] {
	bits := 1 << (indexer.Bits() - 1)

	f := &Filter[T, U]{
		bits: &rye.Bits{
			Data: make(slice.Slice[byte], bits/8),
		},
		idx: indexer,
	}
	f.Add(values...)

	return f
}

// Add values to the Bloom Filter.
func (f *Filter[T, U]) Add(values ...T) error {
	for _, v := range values {
		err := f.add(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Filter[T, U]) add(value T) error {
	idx, err := f.idx.Index(value)
	if err != nil {
		return err
	}

	for _, i := range idx {
		f.bits.Idx = int(i)
		f.bits.Write(1)
	}
	return nil
}

func (f *Filter[T, U]) Contains(value T) bool {
	idx, err := f.idx.Index(value)
	if err != nil {
		return false
	}
	return f.ContainsIndex(idx)
}

func (f *Filter[T, U]) ContainsIndex(idx lhash.Index[U]) bool {
	for _, i := range idx {
		f.bits.Idx = int(i)
		if f.bits.Read() != 1 {
			return false
		}
	}
	return true
}
