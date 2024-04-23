package store

import (
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lstr"
)

var strToBytes = morph.NewValAll(lstr.StringToBytes)

// GetStoresStr is GetStores with the names given as strings.
func GetStoresStr(f FlatFactory, names ...string) (list.Wrapper[FlatStore], error) {
	ns := strToBytes.List(slice.New(names))
	return GetStores(f, ns)
}

// GetStores returns the FlatStore that f holds for each name. The list is lazy:
// f.FlatStore is called every time an index is read, not when GetStores is
// called, so the error returned here is always nil and errors from f are lost.
func GetStores(f FlatFactory, names list.List[[]byte]) (list.Wrapper[FlatStore], error) {
	var errs lerr.Many
	return morph.NewValAll(func(name []byte) FlatStore {
		s, err := f.FlatStore(name)
		errs = errs.Add(err)
		return s
	}).List(names), errs.Cast()
}

// Slice returns every key in the store, in order.
func Slice(s FlatStore) slice.Slice[[]byte] {
	return slice.FromIter(NewIter(s), nil)
}
