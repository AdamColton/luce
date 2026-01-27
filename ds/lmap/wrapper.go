package lmap

// Wrapper provides helpers around a Mapper.
type Wrapper[K comparable, V any] struct {
	Mapper[K, V]
}

// Wrap a Mapper. If it is already a Wrapper, that will be returned.
func Wrap[K comparable, V any](m Mapper[K, V]) Wrapper[K, V] {
	w, ok := m.(Wrapper[K, V])
	if ok {
		return w
	}
	return Wrapper[K, V]{m}
}

// Wrapped fulfills upgrade.Wrapper.
func (w Wrapper[K, V]) Wrapped() any {
	return w.Mapper
}

// GetVal returns the value for a key dropping the "found" boolean.
func (w Wrapper[K, V]) GetVal(key K) (v V) {
	if w.Mapper != nil {
		v, _ = w.Get(key)
	}
	return
}

// Each adds a nil check before calling Mapper.Each.
func (w Wrapper[K, V]) Each(fn EachFunc[K, V]) {
	if w.Mapper == nil {
		return
	}
	w.Mapper.Each(fn)
}

// Len adds a nil check before calling Mapper.Len.
func (w Wrapper[K, V]) Len() int {
	if w.Mapper == nil {
		return 0
	}
	return w.Mapper.Len()
}
