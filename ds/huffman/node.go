package huffman

import (
	"github.com/adamcolton/luce/serial/rye"
)

type huffNode[T any] struct {
	branch [2]*huffNode[T]
	v      T
}

func (d *Frequency[T]) root() *root[T] {
	return newLeaf(d.Val, d.Count)
}

type root[T any] struct {
	node *huffNode[T]
	sum  int
}

func newLeaf[T any](v T, sum int) *root[T] {
	return &root[T]{
		node: &huffNode[T]{
			v: v,
		},
		sum: sum,
	}
}

func newBranch[T any](n0, n1 *huffNode[T], sum int) *root[T] {
	return &root[T]{
		node: &huffNode[T]{
			branch: [2]*huffNode[T]{n0, n1},
		},
		sum: sum,
	}
}

// Read from b until a leaf is encountered, return the leaf value.
func (n *huffNode[T]) Read(b *rye.Bits) T {
	if n.branch[0] == nil {
		return n.v
	}
	return n.branch[b.Read()].Read(b)
}

func (n *huffNode[T]) All(fn func(T)) {
	if n.branch[0] == nil {
		fn(n.v)
	} else {
		n.branch[0].All(fn)
		n.branch[1].All(fn)
	}
}

// func (n *huffNode[T]) gobEncode(enc *gob.Encoder) error {
// 	err := enc.Encode(&(n.v))
// 	if err != nil {
// 		return err
// 	}
// 	if n.branch[0] == nil {
// 		if n.branch[1] != nil {
// 			panic("FAIL (this is a check)")
// 		}
// 		enc.Encode(byte(0))
// 		return nil
// 	}
// 	enc.Encode(byte(1))
// 	err = n.branch[0].gobEncode(enc)
// 	if err != nil {
// 		return err
// 	}
// 	err = n.branch[1].gobEncode(enc)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func decodeHuffNode[T any](dec *gob.Decoder) (*huffNode[T], error) {
// 	n := &huffNode[T]{}
// 	err := dec.Decode(&(n.v))
// 	if err != nil {
// 		return nil, err
// 	}
// 	var b byte
// 	err = dec.Decode(&b)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if b == 1 {
// 		n.branch[0], err = decodeHuffNode[T](dec)
// 		if err != nil {
// 			return nil, err
// 		}
// 		n.branch[1], err = decodeHuffNode[T](dec)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return n, nil
// }
