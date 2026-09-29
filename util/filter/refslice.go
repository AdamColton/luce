package filter

import (
	"github.com/adamcolton/luce/ds/slice"
)

// RefSliceInPlace works like Filter.SliceInPlace, but the filter is called
// with a pointer to each element instead of a copy. This avoids copying large
// values and lets a Filter[*T] be used on a []T. It reorders vals so all the
// elements passing the filter are at the start of the slice and all the
// elements failing the filter are at the end. It returns two subslices, the
// first for passing, the second for failing. No guarantees are made about the
// order of the subslices.
func RefSliceInPlace[T any](f Filter[*T], vals []T) (passing, failing slice.Slice[T]) {
	ln := len(vals)
	if ln == 0 {
		return vals, nil
	}
	start := 0
	end := ln - 1
	for {
		for ; start < ln && f(&vals[start]); start++ {
		}
		for ; end >= 0 && !f(&vals[end]); end-- {
		}
		if start > end {
			break
		}
		vals[start], vals[end] = vals[end], vals[start]
	}
	return vals[:start], vals[start:]
}
