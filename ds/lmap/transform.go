package lmap

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
