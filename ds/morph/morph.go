package morph

type EachFn[K, V any] = func(k K, v V, done *bool)

type Eacher[K, V any] interface {
	Each(fn EachFn[K, V])
}

type KeyVal[K, V, Out any] func(k K, v V) (out Out, include bool)

func NewKeyVal[K, V, Out any](fn KeyVal[K, V, Out]) KeyVal[K, V, Out] {
	return fn
}

type KeyValAll[K, V, Out any] func(k K, v V) (out Out)

func NewKeyValAll[K, V, Out any](fn KeyValAll[K, V, Out]) KeyValAll[K, V, Out] {
	return fn
}

func (kv KeyValAll[K, V, Out]) ToKV() KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		return kv(k, v), true
	}
}

type Val[V, Out any] func(v V) (out Out, include bool)

func NewVal[V, Out any](fn Val[V, Out]) Val[V, Out] {
	return fn
}

type ValAll[V, Out any] func(v V) (out Out)

func NewValAll[V, Out any](fn ValAll[V, Out]) ValAll[V, Out] {
	return fn
}

func (vt ValAll[V, Out]) ToV() Val[V, Out] {
	return func(v V) (out Out, include bool) {
		return vt(v), true
	}
}

type Null[Out any] func() (out Out)

func NewNull[Out any](fn Null[Out]) Null[Out] {
	return fn
}
