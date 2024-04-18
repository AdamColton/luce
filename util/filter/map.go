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
	Each(lmap.EachFunc[K, V])
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

// MapSliceFlag selects what MapFilter.Slice returns, using the flags below.
type MapSliceFlag byte

const (
	// ReturnKeys returns the keys of the pairs that pass the filter.
	ReturnKeys = 1 << iota
	// InverseKeys, with ReturnKeys, returns the keys of the pairs that fail the
	// filter instead.
	InverseKeys
	// ReturnVals returns the values of the pairs that pass the filter.
	ReturnVals
	// InverseVals, with ReturnVals, returns the values of the pairs that fail
	// the filter instead.
	InverseVals

	// ReturnBoth returns both the keys and the values.
	ReturnBoth = ReturnKeys | ReturnVals
)

// Slice runs the filter over every pair in m and returns the keys and/or values
// selected by flags, using keyBuf and valBuf if they have sufficient capacity.
// A slice that flags does not ask for is nil. The order is the order of m.Each.
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
