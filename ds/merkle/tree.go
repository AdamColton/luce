package merkle

import (
	"hash"
	"io"
	"math/bits"

	"github.com/adamcolton/luce/lerr"
)

type tree struct {
	node
	h     hash.Hash
	depth uint32
	pos   int64
}

func (t *tree) updateTree(h hash.Hash, stitchData bool) {
	t.h = h
	ln, _ := t.node.update(h)
	t.depth = uint32(bits.Len(uint(t.Leaves())))
	if stitchData {
		t.stitchData(make([]byte, ln))
	}
}

func (t *tree) Read(p []byte) (n int, err error) {
	if t.pos >= int64(t.Len()) {
		return 0, io.EOF
	}
	n = copy(p, t.Data()[t.pos:])
	t.pos += int64(n)
	return
}

func (t *tree) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		t.pos = offset
	case io.SeekEnd:
		t.pos = int64(t.Len()) + offset
	case io.SeekCurrent:
		t.pos += offset
	default:
		return -1, lerr.ErrBadWhence(whence)
	}

	l64 := int64(t.Len())
	if t.pos < 0 {
		t.pos = 0
	} else if t.pos > l64 {
		t.pos = l64
	}

	return t.pos, nil
}

func (t *tree) Description() Description {
	return Description{
		Digest: t.Digest(),
		Leaves: uint32(t.Leaves()),
	}
}

func (t *tree) Leaf(idx int) *Leaf {
	if idx < 0 || idx >= t.Leaves() {
		return nil
	}
	l := &Leaf{
		Rows: make([]ValidatorRow, 0, t.depth),
	}
	n := t.node
	for {
		if dl, ok := n.(*dataLeaf); ok {
			l.Data = dl.data
			l.Index = dl.idx
			return l
		}
		b := n.(*branch)
		c0ls := b.children[0].Leaves()
		if idx < c0ls {
			n = b.children[0]
			l.Rows = append(l.Rows, ValidatorRow{
				SiblingDigest: b.children[1].Digest(),
				IsFirst:       true,
			})
		} else {
			n = b.children[1]
			idx -= c0ls
			l.Rows = append(l.Rows, ValidatorRow{
				SiblingDigest: b.children[0].Digest(),
				IsFirst:       false,
			})
		}
	}
}
