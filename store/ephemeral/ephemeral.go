package ephemeral

import (
	"fmt"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/store"
)

type record struct {
	*estore
	data []byte
}

type estore struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	records    []*record
}

func (e *estore) Put(key, value []byte) error {
	idx, app := e.idx.Insert(key)
	if !app && e.records[idx] != nil && e.records[idx].estore != nil {
		return fmt.Errorf("Bucket already exists at that key")
	}
	r := &record{
		data: value,
	}
	if app {
		e.records = append(e.records, r)
	} else {
		e.records[idx] = r
	}
	return nil
}

func (e *estore) Get(key []byte) store.Record {
	var out store.Record
	var idx int
	idx, out.Found = e.idx.Get(key)
	if !out.Found {
		return out
	}
	r := e.records[idx]
	out.Store = r.estore
	out.Value = r.data
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

func (e *estore) Store(bkt []byte) (store.Store, error) {
	idx, found := e.idx.Get(bkt)
	if found {
		if e.records[idx].estore == nil {
			return nil, fmt.Errorf("Value already exists at that key")
		}
		return e.records[idx].estore, nil
	}
	r := &record{
		estore: newEstore(e.idxFact, e.bufferSize),
	}
	idx, app := e.idx.Insert(bkt)
	if app {
		e.records = append(e.records, r)
	} else {
		e.records[idx] = r
	}
	return r.estore, nil
}

func newEstore(f byteid.IndexFactory, ln int) *estore {
	return &estore{
		idxFact:    f,
		bufferSize: ln,
		idx:        f(ln),
		records:    make([]*record, ln),
	}
}

type factory struct {
	idxFact    byteid.IndexFactory
	bufferSize int
	idx        byteid.Index
	estores    []*estore
}

func (f *factory) Store(bkt []byte) (store.Store, error) {
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

func Factory(idxFact byteid.IndexFactory, bufferSize int) store.Factory {
	return &factory{
		idxFact:    idxFact,
		idx:        idxFact(bufferSize),
		bufferSize: bufferSize,
		estores:    make([]*estore, bufferSize),
	}
}
