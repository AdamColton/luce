package rbtree

import (
	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/rye/compact"
	lgob "github.com/adamcolton/luce/serial/wrap/gob"
)

func (n *Node[Key, Val]) GobDecode(data []byte) (err error) {
	defer func() {
		lerr.Recover(func(e error) {
			err = e
		})
	}()

	d := compact.NewDeserializer(data)
	n.color = color(d.Byte())

	lgob.Dec(d.CompactSlice(), &(n.KV))
	lgob.Dec(d.CompactSlice(), &(n.chld))
	lgob.Dec(d.CompactSlice(), &(n.prt))

	n.size = int(d.Uint32())
	n.id = d.Uint32()
	return
}

func (n *Node[Key, Val]) GobEncode() (b []byte, err error) {
	defer func() {
		lerr.Recover(func(e error) {
			err = e
		})
	}()

	var ln uint64 = 1 + 4 + 4 // color,size,id

	kv := lgob.Enc(n.KV)
	ln += compact.Size(kv)

	chld := lgob.Enc(n.chld)
	ln += compact.Size(chld)

	prt := lgob.Enc(n.prt)
	ln += compact.Size(prt)

	s := compact.MakeSerializer(int(ln))
	s.Byte(byte(n.color))
	s.CompactSlice(kv)
	s.CompactSlice(chld)
	s.Uint32(uint32(n.size))
	s.Uint32(n.id)
	b = s.Data
	return
}

/*
	root      graph.Ptr[*Node[Key, Val]]
	cmpr      filter.Compare[Key]
	size      int
	idCounter uint32
*/

func (t *Tree[Key, Val]) GobDecode(data []byte) (err error) {
	defer func() {
		lerr.Recover(func(e error) {
			err = e
		})
	}()

	d := compact.NewDeserializer(data)
	var ptr graph.Ptr[*Node[Key, Val]]
	lgob.Dec(d.CompactSlice(), ptr)
	t.root = ptr
	t.size = int(d.Uint32())
	t.idCounter = d.Uint32()

	return
}

func (t *Tree[Key, Val]) GobEncode() (b []byte, err error) {
	defer func() {
		lerr.Recover(func(e error) {
			err = e
		})
	}()

	root := lgob.Enc(t.root)
	ln := compact.Size(root) + 4 + 4 // size, idCounter

	s := compact.MakeSerializer(int(ln))
	s.Make()
	s.CompactSlice(root)
	s.Uint32(uint32(t.size))
	s.Uint32(t.idCounter)
	b = s.Data
	return
}
