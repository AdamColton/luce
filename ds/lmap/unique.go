package lmap

import "github.com/adamcolton/luce/ds/slice"

func Unique[K comparable](s, buf []K) slice.Slice[K] {
	set := make(Map[K, struct{}])
	for _, t := range s {
		set[t] = struct{}{}
	}
	return set.Keys(buf)
}
