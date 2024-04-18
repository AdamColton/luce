package filter

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

// MapFilter provides a filter for Keys and Values. For either, a nil Filter will
// be ignored.
type MapFilter[K comparable, V any] struct {
	Key Filter[K]
	Val Filter[V]
}

// Mapper is a representation of lmap.Mapper. But the filter package only
// needs the Each method.
type Mapper[K comparable, V any] interface {
	Each(lmap.IterFunc[K, V])
	Delete(K)
}

// NewMap creates a Map filter from the provided filters.
func NewMap[K comparable, V any](k Filter[K], v Filter[V]) MapFilter[K, V] {
	return MapFilter[K, V]{
		Key: k,
		Val: v,
	}
}

// Filter checks the key and value against the map filter.
func (mf MapFilter[K, V]) Filter(k K, v V) bool {
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

func (mf MapFilter[K, V]) Slice(m Mapper[K, V], keyBuf []K, valBuf []V, flags MapSliceFlag) (keys slice.Slice[K], vals slice.Slice[V]) {
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

func (mf MapFilter[K, V]) Map(m Mapper[K, V], to lmap.Mapper[K, V]) lmap.Wrapper[K, V] {
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

func (mf MapFilter[K, V]) Purge(m Mapper[K, V], buf []K) slice.Slice[K] {
	rm, _ := mf.Slice(m, buf, nil, ReturnKeys|InverseKeys)
	for _, k := range rm {
		m.Delete(k)
	}
	return rm
}
