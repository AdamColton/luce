package store

import (
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
)

func GetStoresStr(f FlatFactory, names ...string) (list.Wrapper[FlatStore], error) {
	ns := list.TransformSlice(names, func(n string) []byte {
		return []byte(n)
	})
	return GetStores(f, ns)
}

func GetStores(f FlatFactory, names list.List[[]byte]) (list.Wrapper[FlatStore], error) {
	var errs lerr.Many
	return list.NewTransformer(names, func(name []byte) FlatStore {
		s, err := f.FlatStore(name)
		errs = errs.Add(err)
		return s
	}), errs.Cast()
}

func Slice(s FlatStore) slice.Slice[[]byte] {
	return slice.FromIter(NewIter(s), nil)
}
