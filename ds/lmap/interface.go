package lmap

type IterFunc[K comparable, V any] func(key K, val V, done *bool)

type Mapper[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Len() int
	Delete(K)
	Each(IterFunc[K, V])
	Map() map[K]V
}
