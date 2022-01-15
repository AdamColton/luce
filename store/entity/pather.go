package entity

import (
	"reflect"
	"strings"

	"github.com/adamcolton/luce/serial/type32"
)

type EntPathByRelector map[reflect.Type][][]byte

func (ep EntPathByRelector) EntPath(e Entity) ([][]byte, bool) {
	p, ok := ep[reflect.TypeOf(e)]
	return p, ok
}

func NewEntPathByRelector() EntPathByRelector {
	return make(EntPathByRelector)
}

func (ep EntPathByRelector) Add(e Entity, path [][]byte) EntPathByRelector {
	ep[reflect.TypeOf(e)] = path
	return ep
}

func (ep EntPathByRelector) AddString(e Entity, path string) EntPathByRelector {
	return ep.Add(e, StringToPath(path))
}

type EntPather interface {
	EntPath() [][]byte
}

type EntPathByEntPather struct{}

func (ep EntPathByEntPather) EntPath(e Entity) ([][]byte, bool) {
	if ep, ok := e.(EntPather); ok {
		return ep.EntPath(), true
	}
	return nil, false
}

type EntPathByType32 map[uint32][][]byte

func (ep EntPathByType32) EntPath(e Entity) ([][]byte, bool) {
	if t32, ok := e.(type32.TypeIDer32); ok {
		p, ok := ep[t32.TypeID32()]
		return p, ok
	}
	return nil, false
}

type EntPathByMany []Pather

func (ep EntPathByMany) EntPath(e Entity) ([][]byte, bool) {
	for _, epi := range ep {
		p, ok := epi.EntPath(e)
		if ok {
			return p, true
		}
	}
	return nil, false
}

func StringToPath(path string) [][]byte {
	ps := strings.Split(path, "/")
	out := make([][]byte, len(ps))
	for i, p := range ps {
		out[i] = []byte(p)
	}
	return out
}
