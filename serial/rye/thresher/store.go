package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
)

type storeRecord struct {
	data []byte
	t    reflect.Type
}

var store = lmap.Map[string, storeRecord]{}

var encodings = lmap.Map[string, []byte]{}
