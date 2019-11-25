package bytemap

import (
	"bytes"
	"encoding/base64"

	"github.com/adamcolton/luce/ds/idx/byteid"
)

type byteIdxMap struct {
	m        map[string]int
	sliceLen int
	maxIdx   int
	recycle  []int
}

func New(sliceLen int) byteid.Index {
	return &byteIdxMap{
		m:        make(map[string]int, sliceLen),
		sliceLen: sliceLen,
		maxIdx:   0,
	}
}

func (m *byteIdxMap) SliceLen() int {
	return m.sliceLen
}

func (m *byteIdxMap) SetSliceLen(newlen int) {
	if newlen > m.sliceLen {
		m.sliceLen = newlen
	}
}

func (m *byteIdxMap) Insert(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if ok {
		return idx, false
	}
	app := false
	if ln := len(m.recycle); ln > 0 {
		idx = m.recycle[ln-1]
		m.recycle = m.recycle[:ln-1]
	} else {
		idx = m.maxIdx
		m.maxIdx++
		app = m.maxIdx > m.sliceLen
		if app {
			m.sliceLen = m.maxIdx
		}
	}
	m.m[key] = idx
	return idx, app
}

func (m *byteIdxMap) Get(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if !ok {
		return -1, false
	}
	return idx, ok
}

func (m *byteIdxMap) Delete(id []byte) (int, bool) {
	key := base64.StdEncoding.EncodeToString(id)
	idx, ok := m.m[key]
	if !ok {
		return -1, false
	}
	delete(m.m, key)
	m.recycle = append(m.recycle, idx)
	return idx, true
}

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
