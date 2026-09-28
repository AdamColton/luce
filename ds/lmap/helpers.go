package lmap

import "github.com/adamcolton/luce/lerr"

// ErrCollision is returned by CheckSet when the key is already defined.
const ErrCollision = lerr.Str("key collision")

// CheckSet returns ErrCollision if key is already defined, otherwise, it sets
// the given keypair.
func CheckSet[K comparable, V any](m Mapper[K, V], key K, val V) error {
	_, contains := m.Get(key)
	if contains {
		return ErrCollision
	}
	m.Set(key, val)
	return nil
}
