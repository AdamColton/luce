package filter

import "github.com/adamcolton/luce/ds/slice"

type Map[K comparable, V any] struct {
	Key Filter[K]
	Val Filter[V]
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

func (mf Map[K, V]) Slice(m map[K]V, keyBuf []K, valBuf []V, flags MapSliceFlag) (keys slice.Slice[K], vals slice.Slice[V]) {
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
	for k, v := range m {
		f := mf.Filter(k, v)
		if rk && f != ik {
			keys = append(keys, k)
		}
		if rv && f != iv {
			vals = append(vals, v)
		}
	}
	return
}
