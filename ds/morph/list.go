package morph

import "github.com/adamcolton/luce/ds/list"

type List[In, Out any] struct {
	list.List[In]
	ValAll[In, Out]
}

func (l List[In, Out]) AtIdx(i int) Out {
	return l.ValAll(l.List.AtIdx(i))
}

func (va ValAll[In, Out]) List(in list.List[In]) list.Wrapper[Out] {
	return list.Wrap(List[In, Out]{
		List:   in,
		ValAll: va,
	})
}
