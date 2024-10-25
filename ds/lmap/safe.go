package lmap

import "sync"

// Safe provides a map with a Read-Write mutex that provides thread safe access.
type Safe[K comparable, V any] struct {
	m   map[K]V
	mux sync.RWMutex
}

// NewSafe creates a new Wrapped instance of a Safe map.
func NewSafe[K comparable, V any](m map[K]V) Wrapper[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return Wrap(&Safe[K, V]{
		m: m,
	})
}

// EmptySafe creates a new empty Wrapped instance of a Safe map with the defined
// capacity.
func EmptySafe[K comparable, V any](capacity int) Wrapper[K, V] {
	return Wrap(&Safe[K, V]{
		m: make(map[K]V, capacity),
	})
}

// Len fulfills Mapper, returning the length of the underlying map.
func (s *Safe[K, V]) Len() int {
	return len(s.m)
}

// Get fulfills Mapper, returning the value for a given key and a boolean
// indicating if the key was present.
func (s *Safe[K, V]) Get(key K) (V, bool) {
	s.mux.RLock()
	v, ok := s.m[key]
	s.mux.RUnlock()
	return v, ok
}

// Set fulfills Mapper, setting the key and value in the underlying map.
func (s *Safe[K, V]) Set(key K, val V) {
	s.mux.Lock()
	s.m[key] = val
	s.mux.Unlock()
}

// Delete fulfills Mapper, deleting the key from the underlying map.
func (s *Safe[K, V]) Delete(key K) {
	s.mux.Lock()
	delete(s.m, key)
	s.mux.Unlock()
}

// Each calls fn for every key-value pair.
func (s *Safe[K, V]) Each(fn IterFunc[K, V]) {
	done := false
	s.mux.RLock()
	for k, v := range s.m {
		fn(k, v, &done)
		if done {
			break
		}
	}
	s.mux.RUnlock()
}

// Map fulfills Mapper and returns the underlying map.
func (s *Safe[K, V]) Map() map[K]V {
	return s.m
}
