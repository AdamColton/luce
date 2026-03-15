package morph

import "github.com/adamcolton/luce/ds/list"

// List is a list.List that applies a ValAll to the values of another List when
// they are read. Nothing is cached, so the ValAll is called on every AtIdx.
type List[In, Out any] struct {
	list.List[In]
	ValAll[In, Out]
}

// AtIdx fulfills list.List. It applies the ValAll to the value at i of the
// underlying List.
func (l List[In, Out]) AtIdx(i int) Out {
	return l.ValAll(l.List.AtIdx(i))
}

// List returns a list.Wrapper of the results of applying the ValAll to the
// values of the List. The Len is that of the List.
func (va ValAll[In, Out]) List(in list.List[In]) list.Wrapper[Out] {
	return list.Wrap(List[In, Out]{
		List:   in,
		ValAll: va,
	})
}
