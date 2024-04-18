package lmap

type Map[K comparable, V any] map[K]V

func (m Map[K, V]) Len() int {
	return len(m)
}
