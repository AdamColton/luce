package lset

import "github.com/adamcolton/luce/ds/slice"

// Unique is a convenience function that casts a slice to a Set and back to
// make it unique.
func Unique[K comparable](s, buf []K) slice.Slice[K] {
	return New(s...).Slice(buf)
}
