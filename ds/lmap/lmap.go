package lmap

// Map fulfills Mapper using the builtin map type.
type Map[K comparable, V any] map[K]V

// New creates a Wrapped instance of Map.
func New[K comparable, V any](m map[K]V) Wrapper[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return Wrap(Map[K, V](m))
}

// Empty creates a Wrapped instance of Map with the defined capacity.
func Empty[K comparable, V any](capacity int) Wrapper[K, V] {
	return Wrap(make(Map[K, V], capacity))
}

// Len fulfills Mapper, returning the length of the underlying map.
func (m Map[K, V]) Len() int {
	return len(m)
}

// Get fulfills Mapper, returning the value for a given key and a boolean
// indicating if the key was present.
func (m Map[K, V]) Get(key K) (V, bool) {
	v, ok := m[key]
	return v, ok
}

// Set fulfills Mapper, setting the key and value in the underlying map.
func (m Map[K, V]) Set(key K, val V) {
	m[key] = val
}

// Delete fulfills Mapper, deleting the key from the underlying map.
func (m Map[K, V]) Delete(key K) {
	delete(m, key)
}

// Each calls fn for every key-value pair.
func (m Map[K, V]) Each(fn IterFunc[K, V]) {
	done := false
	for k, v := range m {
		fn(k, v, &done)
		if done {
			break
		}
	}
}

// Map fulfills Mapper and returns the underlying map.
func (m Map[K, V]) Map() map[K]V {
	return m
}

// New fulfills Mapper, returning a new Map
func (m Map[K, V]) New() Mapper[K, V] {
	return make(Map[K, V], 0)
}
