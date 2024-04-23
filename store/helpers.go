package store

import (
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/lerr"
)

func GetStoresStr(f Factory, names ...string) (list.Wrapper[Store], error) {
	ns := list.TransformSlice(names, func(n string) []byte {
		return []byte(n)
	})
	return GetStores(f, ns)
}

func GetStores(f Factory, names list.List[[]byte]) (list.Wrapper[Store], error) {
	var errs lerr.Many
	return list.NewTransformer(names, func(name []byte) Store {
		s, err := f.Store(name)
		errs = errs.Add(err)
		return s
	}), errs.Cast()
}
