package bytemap

import (
	"bytes"
	"encoding/base64"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/adamcolton/luce/ds/idx/sliceidx"
)

type byteIdxMap struct {
	m  map[string]int
	si sliceidx.SliceIdx
}

// New fulfills byteid.Factory. It returns an instance of byteid. Index that is
// backed by using a map to map the id to an int. The key is formed by
// converting the id to base64 encoded string.
func New(sliceLen int) byteid.Index {
	return &byteIdxMap{
		m:  make(map[string]int, sliceLen),
		si: sliceidx.New(sliceLen),
	}
}

// SliceLen of the Indexed slice.
func (m *byteIdxMap) SliceLen() int {
	return m.si.SliceLen
}

// SetSliceLen can be used to grow the slice.
func (m *byteIdxMap) SetSliceLen(newlen int) {
	m.si.SetSliceLen(newlen)
}

// Insert an ID. The first value returned is the index and the bool
// indicates if an append is required.
func (m *byteIdxMap) Insert(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if ok {
		return idx, false
	}
	idx, app := m.si.NextIdx()
	m.m[key] = idx
	return idx, app
}

// Get by ID. If not found it will return (-1,false). If it is found the
// first value is the index and the second value is True.
func (m *byteIdxMap) Get(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if !ok {
		return -1, false
	}
	return idx, ok
}

// Delete by ID. Removes the ID from the index, the value will be recycled. This
// should be called before removing the value from the slice.
func (m *byteIdxMap) Delete(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if !ok {
		return -1, false
	}
	delete(m.m, key)
	m.si.Recycle(idx)
	return idx, true
}

// Next ID after the ID given
func (m *byteIdxMap) Next(after []byte) ([]byte, int) {
	var best []byte
	bestIdx := -1
	for idStr, idx := range m.m {
		id, _ := base64.StdEncoding.DecodeString(idStr)
		a := bytes.Compare(after, id)
		if a == -1 && (bestIdx == -1 || bytes.Compare(id, best) == -1) {
			best = id
			bestIdx = idx
		}
	}
	return best, bestIdx
}
