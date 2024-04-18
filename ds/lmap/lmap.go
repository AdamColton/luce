package lmap

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) Pop(key K) (V, bool) {
	v, found := m[key]
	if found {
		delete(m, key)
	}
	return v, found
}
