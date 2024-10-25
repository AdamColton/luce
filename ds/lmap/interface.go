package lmap

// IterFunc is a function that can be called by Each on a Mapper. Note that
// "done" is a return argument. The choice was made in this case because
// being able to stop the iteration is useful in enough cases to justify it's
// inclusion, it is used infrequently. Not requiring a return argumenet cleaned
// up most of the instances of IterFuncs.
type IterFunc[K comparable, V any] func(key K, val V, done *bool)

// Mapper represents the operations a Map can perform.
type Mapper[K comparable, V any] interface {
	Get(K) (V, bool)
	Set(K, V)
	Len() int
	Delete(K)
	Each(IterFunc[K, V])
	Map() map[K]V
}
