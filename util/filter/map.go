package filter

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

type Map[K comparable, V any] struct {
	Key Filter[K]
	Val Filter[V]
}

type Mapper[K comparable, V any] interface {
	Each(lmap.IterFunc[K, V])
}

func NewMap[K comparable, V any](k Filter[K], v Filter[V]) Map[K, V] {
	return Map[K, V]{
		Key: k,
		Val: v,
	}
}

func (mf Map[K, V]) Filter(k K, v V) bool {
	return (mf.Key == nil || mf.Key(k)) && (mf.Val == nil || mf.Val(v))
}

type MapSliceFlag byte

const (
	ReturnKeys = 1 << iota
	InverseKeys
	ReturnVals
	InverseVals

	ReturnBoth = ReturnKeys | ReturnVals
)

func (mf Map[K, V]) Slice(m Mapper[K, V], keyBuf []K, valBuf []V, flags MapSliceFlag) (keys slice.Slice[K], vals slice.Slice[V]) {
	rk := flags&ReturnKeys == ReturnKeys
	ik := flags&InverseKeys == InverseKeys
	rv := flags&ReturnVals == ReturnVals
	iv := flags&InverseVals == InverseVals

	if !rk && !rv {
		return
	}
	if rk {
		keys = keyBuf[:0]
	}
	if rv {
		vals = valBuf[:0]
	}
	m.Each(func(key K, val V, done *bool) {
		f := mf.Filter(key, val)
		if rk && f != ik {
			keys = append(keys, key)
		}
		if rv && f != iv {
			vals = append(vals, val)
		}
	})
	return
}

func (mf Map[K, V]) Map(m Mapper[K, V], to lmap.Mapper[K, V]) lmap.Wrapper[K, V] {
	if to == nil {
		to = lmap.New[K, V](nil)
	}
	m.Each(func(key K, val V, done *bool) {
		if mf.Filter(key, val) {
			to.Set(key, val)
		}
	})
	return lmap.Wrap(to)
}
