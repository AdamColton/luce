package ll

import (
	"github.com/adamcolton/luce/ds/graph"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/liter"
)

const ErrPtrMiss = lerr.Str("bad graph.Ptr")

type Single[Key, Val any] struct {
	graph.KV[Key, Val]
	Next graph.Ptr[*Single[Key, Val]]
}

func NewSingle[Key, Val any](ptr graph.Ptr[*Single[Key, Val]], k Key, v Val) *Single[Key, Val] {
	return &Single[Key, Val]{
		KV:   graph.NewKV(k, v),
		Next: ptr,
	}
}

func NewSingleLoop[Key, Val any](ptr graph.Ptr[*Single[Key, Val]], k Key, v Val) *Single[Key, Val] {
	n := NewSingle(ptr, k, v)
	n.Next = n.Next.Set(n)
	return n
}

func (s *Single[Key, Val]) InsertAfter(k Key, v Val) *Single[Key, Val] {
	n := &Single[Key, Val]{
		KV:   graph.NewKV(k, v),
		Next: s.Next,
	}
	s.Next = s.Next.New().Set(n)
	return n
}

func (s *Single[Key, Val]) DeleteNext() {
	nxt := lerr.OK(s.Next.Get())(ErrPtrMiss)
	if nxt == nil {
		return
	}
	s.Next = nxt.Next
}

func (s Single[Key, Val]) AtIdx(idx int) graph.Node[Key, Val] {
	if idx == 0 {
		return lerr.OK(s.Next.Get())(ErrPtrMiss)
	}
	return nil
}

func (s *Single[Key, Val]) Len() int {
	return 1
}

func (s *Single[Key, Val]) Iter() liter.Wrapper[graph.KV[Key, Val]] {
	return liter.Wrap(&singleIter[Key, Val]{
		n:     s,
		start: s,
	})
}

func (s *Single[Key, Val]) IterFactory() (iter liter.Iter[graph.KV[Key, Val]], v Val, done bool) {
	done = s == nil
	if !done {
		v = s.V
	}
	iter = s.Iter()
	return
}

type singleIter[Key, Val any] struct {
	i     int
	n     *Single[Key, Val]
	start *Single[Key, Val]
}

func (si *singleIter[Key, Val]) Next() (kv graph.KV[Key, Val], done bool) {
	done = si.Done()
	if done {
		return
	}
	si.i++
	si.n = lerr.OK(si.n.Next.Get())(ErrPtrMiss)
	done = si.Done()
	if !done {
		kv = si.n.KV
	}
	return
}

func (si *singleIter[Key, Val]) Cur() (kv graph.KV[Key, Val], done bool) {
	done = si.Done()
	if !done {
		kv = si.n.KV
	}
	return
}

func (si *singleIter[Key, Val]) Done() bool {
	return si.n == nil || (si.i > 0 && si.n == si.start)
}

func (si *singleIter[Key, Val]) Idx() int {
	return si.i
}
