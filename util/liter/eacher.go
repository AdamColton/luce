package liter

type Eacher[V any] struct {
	Func func(inner EachFn[V])
	L    int
}

func (e Eacher[V]) Each(fn EachFn[V]) {
	e.Func(fn)
}

func (e Eacher[V]) Len() int {
	return e.L
}
