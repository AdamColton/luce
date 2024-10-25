package lmap

import "sync"

type Safe[K comparable, V any] struct {
	m   map[K]V
	mux sync.RWMutex
}

func NewSafe[K comparable, V any](m map[K]V) Wrapper[K, V] {
	if m == nil {
		m = make(map[K]V)
	}
	return Wrap(&Safe[K, V]{
		m: m,
	})
}

func EmptySafe[K comparable, V any](ln int) Wrapper[K, V] {
	return Wrap(&Safe[K, V]{
		m: make(map[K]V, ln),
	})
}

func (s *Safe[K, V]) Len() int {
	return len(s.m)
}

func (s *Safe[K, V]) Get(key K) (V, bool) {
	s.mux.RLock()
	v, ok := s.m[key]
	s.mux.RUnlock()
	return v, ok
}

func (s *Safe[K, V]) Set(key K, val V) {
	s.mux.Lock()
	s.m[key] = val
	s.mux.Unlock()
}

func (s *Safe[K, V]) Delete(key K) {
	s.mux.Lock()
	delete(s.m, key)
	s.mux.Unlock()
}

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

func (s *Safe[K, V]) Map() map[K]V {
	return s.m
}
