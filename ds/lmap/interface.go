package lmap

// EachFunc is a function that can be called by Each on a Mapper. Note that
// "done" is a return argument. The choice was made in this case because
// being able to stop the iteration is useful in enough cases to justify its
// inclusion, it is used infrequently. Not requiring a return argument cleaned
// up most of the instances of EachFunc.
type EachFunc[K comparable, V any] = func(key K, val V, done *bool)

// Eacher is fulfilled by any type that can call a function on each of its
// key/value pairs. Setting done to true in the function stops the iteration.
type Eacher[K comparable, V any] interface {
	Each(EachFunc[K, V])
}

// [ ] Set(K, V) V
//     have set return the value, better for chaining

// Mapper represents the operations a Map can perform.
type Mapper[K comparable, V any] interface {
	MapReader[K, V]
	// Set sets the value for the key, replacing any existing value.
	Set(K, V)
	// Delete removes the key. It is not an error if the key is absent.
	Delete(K)
}

// MapReader is the read-only part of Mapper.
type MapReader[K comparable, V any] interface {
	// Get returns the value for the key and true if the key is present.
	Get(K) (V, bool)
	// Len returns the number of key/value pairs.
	Len() int
	Eacher[K, V]
	// Map returns the underlying map, not a copy.
	Map() map[K]V
	// New creates a new Mapper with the same underlying structure.
	New() Mapper[K, V]
}
