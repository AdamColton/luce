package store

import (
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lstr"
)

var strToBytes = list.TransformAny(lstr.StringToBytes)

func GetStoresStr(f FlatFactory, names ...string) (list.Wrapper[FlatStore], error) {
	ns := strToBytes.New(slice.New(names))
	return GetStores(f, ns)
}

func GetStores(f FlatFactory, names list.List[[]byte]) (list.Wrapper[FlatStore], error) {
	var errs lerr.Many
	return list.NewTransformer(names, func(name []byte, idx int) FlatStore {
		s, err := f.FlatStore(name)
		errs = errs.Add(err)
		return s
	}), errs.Cast()
}

func Slice(s FlatStore) slice.Slice[[]byte] {
	return slice.FromIter(NewIter(s), nil)
}
