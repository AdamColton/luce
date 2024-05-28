package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/bimap"
)

var (
	typeIDs        = bimap.New[uint32, reflect.Type](0)
	maxType uint32 = 1
)

func type2id(t reflect.Type) uint32 {
	id, found := typeIDs.B(t)
	if !found {
		id = maxType
		maxType++
		typeIDs.Add(id, t)
	}
	return id
}
