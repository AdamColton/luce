package lmap

// EachFunc is a function that can be called by Each on a Mapper. Note that
// "done" is a return argument. The choice was made in this case because
// being able to stop the iteration is useful in enough cases to justify it's
// inclusion, it is used infrequently. Not requiring a return argumenet cleaned
// up most of the instances of EachFunc.
type EachFunc[K comparable, V any] = func(key K, val V, done *bool)

type Eacher[K comparable, V any] interface {
	Each(EachFunc[K, V])
}

// Mapper represents the operations a Map can perform.
type Mapper[K comparable, V any] interface {
	MapReader[K, V]
	Set(K, V)
	Delete(K)
}

type MapReader[K comparable, V any] interface {
	Get(K) (V, bool)
	Len() int
	Eacher[K, V]
	Map() map[K]V
	// New creates a new Mapper with the same underlying structure.
	New() Mapper[K, V]
}
