package lmap

import "github.com/adamcolton/luce/ds/slice"

// ForAll is a helper function. The output will generally be fed into either
// TransformVal or TransformKey.
func ForAll[In, Out any](fn func(In) Out) func(In) (Out, bool) {
	return func(in In) (Out, bool) {
		return fn(in), true
	}
}

// TransformVal is a helper that applies the given function to the value.
func TransformVal[Key comparable, VIn any, VOut any](fn func(v VIn) (VOut, bool)) TransformFunc[Key, VIn, Key, VOut] {
	return func(k Key, v VIn) (Key, VOut, bool) {
		vo, ok := fn(v)
		return k, vo, ok
	}
}

// TransformKey is a helper that applies the given function to the key.
func TransformKey[V any, KIn, KOut comparable](fn func(k KIn) (KOut, bool)) TransformFunc[KIn, V, KOut, V] {
	return func(k KIn, v V) (KOut, V, bool) {
		ko, ok := fn(k)
		return ko, v, ok
	}
}

// TransformFunc converts the key and value types. Only key/values pairs for
// which include is true are used.
type TransformFunc[KIn comparable, VIn any, KOut comparable, VOut any] func(k KIn, v VIn) (kOut KOut, vOut VOut, include bool)

// Transform applies the TransformFunc to m and sets the results on buf
// when include is true. Note that buf is not cleared, so this can perform as
// and append. If buf is nil a lmap.Map is created sized to m.
func (fn TransformFunc[KIn, VIn, KOut, VOut]) Transform(m Mapper[KIn, VIn], buf Mapper[KOut, VOut]) Wrapper[KOut, VOut] {
	var out Wrapper[KOut, VOut]
	if buf == nil {
		out = Empty[KOut, VOut](m.Len())
	} else {
		out = Wrap(buf)
	}
	m.Each(func(key KIn, val VIn, done *bool) {
		ko, vo, ok := fn(key, val)
		if ok {
			out.Set(ko, vo)
		}
	})
	return out
}

// Map takes in a plain map and creates an instance of lmap.Map.
func (fn TransformFunc[KIn, VIn, KOut, VOut]) Map(m map[KIn]VIn) Wrapper[KOut, VOut] {
	return fn.Transform(New(m), nil)
}

// NewTransformFunc is just a helper for converting a function to the the
// TransformFunc type.
func NewTransformFunc[KIn comparable, VIn any, KOut comparable, VOut any](fn func(k KIn, v VIn) (KOut, VOut, bool)) TransformFunc[KIn, VIn, KOut, VOut] {
	return fn
}

// TransformMap takes in a plain map and uses the TransformFunc to create an
// instance of lmap.Map.
func TransformMap[KIn comparable, VIn any, KOut comparable, VOut any](m map[KIn]VIn, fn TransformFunc[KIn, VIn, KOut, VOut]) Wrapper[KOut, VOut] {
	return fn.Map(m)
}

// Transform applies the TransformFunc to m and sets the results on buf when
// include is true. Note that buf is not cleared, so this can perform as and
// append. If buf is nil a lmap.Map is created sized to m.
func Transform[KIn comparable, VIn any, KOut comparable, VOut any](m Mapper[KIn, VIn], buf Mapper[KOut, VOut], fn TransformFunc[KIn, VIn, KOut, VOut]) Wrapper[KOut, VOut] {
	return fn.Transform(m, buf)
}

// SliceTransformFunc is used to transform the values in a map to slice.
type SliceTransformFunc[K comparable, V any, Out any] func(k K, v V) (out Out, include bool)

// NewSliceTransformFunc is just a helper for converting a function to the the
// SliceTransformFunc type.
func NewSliceTransformFunc[K comparable, V any, Out any](fn func(k K, v V) (out Out, include bool)) SliceTransformFunc[K, V, Out] {
	return fn
}

// Transform uses the SliceTransformFunc to transform the map to a slice.
func (fn SliceTransformFunc[K, V, Out]) Transform(m Mapper[K, V], buf []Out) slice.Slice[Out] {
	if buf == nil {
		buf = make([]Out, 0, m.Len())
	}
	m.Each(func(key K, val V, done *bool) {
		v, ok := fn(key, val)
		if ok {
			buf = append(buf, v)
		}
	})
	return buf
}

// TransformMap uses the SliceTransformFunc to transform the map to a slice.
func (fn SliceTransformFunc[K, V, Out]) TransformMap(m map[K]V) slice.Slice[Out] {
	return fn.Transform(New(m), nil)
}

// SliceTransformMap takes a function for producing a single value from a given
// key and value produces a map.
func SliceTransformMap[K comparable, V any, Out any](m map[K]V, fn SliceTransformFunc[K, V, Out]) slice.Slice[Out] {
	return fn.Transform(New(m), nil)
}

// SliceTransform takes a function for producing a single value from a given key
// and value produces a map.
func SliceTransform[K comparable, V any, Out any](m Mapper[K, V], buf []Out, fn SliceTransformFunc[K, V, Out]) slice.Slice[Out] {
	return fn.Transform(m, buf)
}
