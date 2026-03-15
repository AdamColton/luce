package morph

// EachFn is the callback used by Eacher. It is called once for each key/value
// pair, and setting done to true stops the iteration. It is the same as
// lmap.EachFunc.
type EachFn[K, V any] = func(k K, v V, done *bool)

// Eacher is fulfilled by any type that can call a function on each of its
// key/value pairs.
type Eacher[K, V any] interface {
	Each(fn EachFn[K, V])
}

// KeyVal converts a key and a value to an Out. The tools that apply a KeyVal
// skip the pairs for which include is false.
type KeyVal[K, V, Out any] func(k K, v V) (out Out, include bool)

// NewKeyVal is syntactic sugar to infer the types of a KeyVal.
func NewKeyVal[K, V, Out any](fn KeyVal[K, V, Out]) KeyVal[K, V, Out] {
	return fn
}

// KeyValAll converts a key and a value to an Out and includes every pair.
type KeyValAll[K, V, Out any] func(k K, v V) (out Out)

// NewKeyValAll is syntactic sugar to infer the types of a KeyValAll.
func NewKeyValAll[K, V, Out any](fn KeyValAll[K, V, Out]) KeyValAll[K, V, Out] {
	return fn
}

// ToKV converts the KeyValAll to a KeyVal that always includes the pair.
func (kv KeyValAll[K, V, Out]) ToKV() KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		return kv(k, v), true
	}
}

// GetKey returns a KeyValAll that returns the key and ignores the value.
func GetKey[K, V any]() KeyValAll[K, V, K] {
	return func(k K, v V) K {
		return k
	}
}

// Val converts a value to an Out. The tools that apply a Val skip the values for
// which include is false.
type Val[V, Out any] func(v V) (out Out, include bool)

// NewVal is syntactic sugar to infer the types of a Val.
func NewVal[V, Out any](fn Val[V, Out]) Val[V, Out] {
	return fn
}

// ValAll converts a value to an Out and includes every value.
type ValAll[V, Out any] func(v V) (out Out)

// NewValAll is syntactic sugar to infer the types of a ValAll.
func NewValAll[V, Out any](fn ValAll[V, Out]) ValAll[V, Out] {
	return fn
}

// ToV converts the ValAll to a Val that always includes the value.
func (vt ValAll[V, Out]) ToV() Val[V, Out] {
	return func(v V) (out Out, include bool) {
		return vt(v), true
	}
}

// Null creates an Out without any input.
type Null[Out any] func() (out Out)

// NewNull is syntactic sugar to infer the type of a Null.
func NewNull[Out any](fn Null[Out]) Null[Out] {
	return fn
}
