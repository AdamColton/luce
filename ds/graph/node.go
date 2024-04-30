package graph

import "github.com/adamcolton/luce/ds/list"

type Node[Key, Val any] interface {
	list.List[Node[Key, Val]]
	Val() Val
	Key() Key
}

type KV[Key, Val any] struct {
	K Key
	V Val
}

func NewKV[Key, Val any](k Key, v Val) KV[Key, Val] {
	return KV[Key, Val]{
		K: k,
		V: v,
	}
}

func (kv KV[Key, Val]) Key() Key {
	return kv.K
}

func (kv KV[Key, Val]) Val() Val {
	return kv.V
}

func (kv KV[Key, Val]) KeyVal() KV[Key, Val] {
	return kv
}

type Ptr[Val any] interface {
	Get() Val
	Set(Val) Ptr[Val]
	New() Ptr[Val]
}

// I need a way to set this relativly. For instance, if Ptr is holding file
// offsets, I should be able to set with a new offset

type RawPointer[Node any] struct {
	v *Node
}

func (p RawPointer[Node]) Get() *Node {
	return p.v
}

func (p RawPointer[Node]) Set(v *Node) Ptr[*Node] {
	p.v = v
	return p
}

func (p RawPointer[Node]) New() Ptr[*Node] {
	return RawPointer[Node]{}
}
