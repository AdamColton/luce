package lmap

type Map[K comparable, V any] map[K]V

func New[K comparable, V any](m map[K]V) Wrapper[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return Wrap(Map[K, V](m))
}

func Empty[K comparable, V any](ln int) Wrapper[K, V] {
	return Wrap(make(Map[K, V], ln))
}

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

func (m Map[K, V]) Set(key K, val V) {
	m[key] = val
}

func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

func (m Map[K, V]) Each(fn IterFunc[K, V]) {
	done := false
	for k, v := range m {
		fn(k, v, &done)
		if done {
			break
		}
	}
}

func (m Map[K, V]) Map() map[K]V {
	return m
}
