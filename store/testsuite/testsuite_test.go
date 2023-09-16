package testsuite_test

import (
	"errors"
	"sort"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/testsuite"
)

// mem is the simplest store that passes the suite: two maps, one for values and
// one for Sub-Stores.
type mem struct {
	values map[string][]byte
	stores map[string]*mem
}

func newMem() *mem {
	return &mem{
		values: make(map[string][]byte),
		stores: make(map[string]*mem),
	}
}

func (m *mem) Put(key, value []byte) error {
	if _, found := m.stores[string(key)]; found {
		return errors.New("a Sub-Store is at that key")
	}
	m.values[string(key)] = append([]byte(nil), value...)
	return nil
}

func (m *mem) Get(key []byte) store.Record {
	if s, found := m.stores[string(key)]; found {
		return store.Record{Found: true, Store: s}
	}
	v, found := m.values[string(key)]
	return store.Record{Found: found, Value: v}
}

func (m *mem) Next(key []byte) []byte {
	var keys []string
	for k := range m.values {
		keys = append(keys, k)
	}
	for k := range m.stores {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if key == nil || k > string(key) {
			return []byte(k)
		}
	}
	return nil
}

func (m *mem) Delete(key []byte) error {
	delete(m.values, string(key))
	delete(m.stores, string(key))
	return nil
}

func (m *mem) Len() int {
	return len(m.values) + len(m.stores)
}

func (m *mem) NestedStore(name []byte) (store.NestedStore, error) {
	if _, found := m.values[string(name)]; found {
		return nil, errors.New("a value is at that key")
	}
	s, found := m.stores[string(name)]
	if !found {
		s = newMem()
		m.stores[string(name)] = s
	}
	return s, nil
}

func TestAll(t *testing.T) {
	testsuite.TestAll(t, newMem())
}
