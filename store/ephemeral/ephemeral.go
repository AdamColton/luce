package ephemeral

import (
	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/store"
)

type estore struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	records    [][]byte
}

func (e *estore) Len() int {
	return e.idx.Len()
}

func (e *estore) Put(key, value []byte) error {
	idx, app := e.idx.Insert(key)
	if app {
		e.records = append(e.records, value)
	} else {
		e.records[idx] = value
	}
	return nil
}

func (e *estore) Get(key []byte) store.Record {
	var out store.Record
	var idx int
	idx, out.Found = e.idx.Get(key)
	if out.Found {
		out.Value = e.records[idx]
	}
	return out
}

func (e *estore) Next(key []byte) []byte {
	next, _ := e.idx.Next(key)
	return next
}

func (e *estore) Delete(key []byte) error {
	idx, found := e.idx.Delete(key)
	if found {
		e.records[idx] = nil
	}
	return nil
}

func newEstore(f byteid.IndexFactory, ln int) *estore {
	return &estore{
		idxFact:    f,
		bufferSize: ln,
		idx:        f(ln),
		records:    make([][]byte, ln),
	}
}

type factory struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	estores    []*estore
}

func (f *factory) FlatStore(bkt []byte) (store.FlatStore, error) {
	idx, found := f.idx.Get(bkt)
	if found {
		return f.estores[idx], nil
	}
	e := newEstore(f.idxFact, f.bufferSize)
	idx, app := f.idx.Insert(bkt)
	if app {
		f.estores = append(f.estores, e)
	} else {
		f.estores[idx] = e
	}
	return e, nil
}

func Factory(idxFact byteid.IndexFactory, bufferSize int) store.FlatFactory {
	return &factory{
		idxFact:    idxFact,
		idx:        idxFact(bufferSize),
		bufferSize: bufferSize,
		estores:    make([]*estore, bufferSize),
	}
}
