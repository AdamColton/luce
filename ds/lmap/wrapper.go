package lmap

type Wrapper[K comparable, V any] struct {
	Mapper[K, V]
}

func Wrap[K comparable, V any](m Mapper[K, V]) Wrapper[K, V] {
	w, ok := m.(Wrapper[K, V])
	if ok {
		return w
	}
	return Wrapper[K, V]{m}
}

func (w Wrapper[K, V]) Wrapped() any {
	return w.Mapper
}

func (w Wrapper[K, V]) GetVal(key K) V {
	t, _ := w.Get(key)
	return t
}
