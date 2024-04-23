package store

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
)

func GetStoresStr(f Factory, names ...string) (slice.Slice[Store], error) {
	ns := slice.TransformSlice(names, func(n string, _ int) ([]byte, bool) {
		return []byte(n), true
	})
	return GetStores(f, ns...)
}

func GetStores(f Factory, names ...[]byte) (slice.Slice[Store], error) {
	var errs lerr.Many
	return slice.TransformSlice(names, func(name []byte, _ int) (Store, bool) {
		s, err := f.Store(name)
		errs = errs.Add(err)
		return s, err == nil
	}), errs.Cast()
}

func Slice(s Store) slice.Slice[[]byte] {
	return slice.FromIter(NewIter(s), nil)
}
