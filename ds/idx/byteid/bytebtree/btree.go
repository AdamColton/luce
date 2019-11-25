package bytebtree

import (
	"bytes"

	"github.com/google/btree"
)

type entry struct {
	id  []byte
	idx int
}

func (e entry) Less(i btree.Item) bool {
	var id []byte
	switch i := i.(type) {
	case entry:
		id = i.id
	case wrap:
		id = i
	}
	return bytes.Compare(e.id, id) < 0
}

type wrap []byte

func (id wrap) Less(i btree.Item) bool {
	return bytes.Compare(id, i.(entry).id) == -1
}
