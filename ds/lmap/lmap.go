package lmap

type Map[K comparable, V any] map[K]V

func New[K comparable, V any](m map[K]V) Map[K, V] {
	return Map[K, V](m)
}

func (m Map[K, V]) Len() int {
	return len(m)
}
