package morph

import "github.com/adamcolton/luce/util/filter"

// FilterKey converts the KeyValAll to a KeyVal that only includes the pairs for
// which the filter is true for the key. The KeyValAll is only called for those.
func (kvt KeyValAll[K, V, Out]) FilterKey(f filter.Filter[K]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(k)
		if include {
			out = kvt(k, v)
		}
		return
	}
}

// FilterVal converts the KeyValAll to a KeyVal that only includes the pairs for
// which the filter is true for the value. The KeyValAll is only called for
// those.
func (kvt KeyValAll[K, V, Out]) FilterVal(f filter.Filter[V]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(v)
		if include {
			out = kvt(k, v)
		}
		return
	}
}

// FilterOut converts the KeyValAll to a KeyVal that only includes the pairs for
// which the filter is true for the result. The KeyValAll is called for every
// pair.
func (kvt KeyValAll[K, V, Out]) FilterOut(f filter.Filter[Out]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		out = kvt(k, v)
		include = f(out)
		return
	}
}

// FilterVal converts the ValAll to a Val that only includes the values for
// which the filter is true. The ValAll is only called for those.
func (vt ValAll[V, Out]) FilterVal(f filter.Filter[V]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		include = f(v)
		if include {
			out = vt(v)
		}
		return
	}
}

// FilterOut converts the ValAll to a Val that only includes the values for
// which the filter is true for the result. The ValAll is called for every
// value.
func (vt ValAll[V, Out]) FilterOut(f filter.Filter[Out]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		out = vt(v)
		include = f(out)
		return
	}
}

// FilterKey returns a KeyVal that also requires the filter to be true for the
// key. The KeyVal is only called if it is.
func (kvt KeyVal[K, V, Out]) FilterKey(f filter.Filter[K]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(k)
		if include {
			out, include = kvt(k, v)
		}
		return
	}
}

// FilterVal returns a KeyVal that also requires the filter to be true for the
// value. The KeyVal is only called if it is.
func (kvt KeyVal[K, V, Out]) FilterVal(f filter.Filter[V]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(v)
		if include {
			out, include = kvt(k, v)
		}
		return
	}
}

// FilterOut returns a KeyVal that also requires the filter to be true for the
// result. The filter is only called if the KeyVal includes the pair.
func (kvt KeyVal[K, V, Out]) FilterOut(f filter.Filter[Out]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		out, include = kvt(k, v)
		if include {
			include = f(out)
		}
		return
	}
}

// FilterVal returns a Val that also requires the filter to be true for the
// value. The Val is only called if it is.
func (vt Val[V, Out]) FilterVal(f filter.Filter[V]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		include = f(v)
		if include {
			out, include = vt(v)
		}
		return
	}
}

// FilterOut returns a Val that also requires the filter to be true for the
// result. The filter is only called if the Val includes the value.
func (vt Val[V, Out]) FilterOut(f filter.Filter[Out]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		out, include = vt(v)
		if include {
			include = f(out)
		}
		return
	}
}
