package morph

import "github.com/adamcolton/luce/util/filter"

func (kvt KeyValAll[K, V, Out]) FilterKey(f filter.Filter[K]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(k)
		if include {
			out = kvt(k, v)
		}
		return
	}
}

func (kvt KeyValAll[K, V, Out]) FilterVal(f filter.Filter[V]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(v)
		if include {
			out = kvt(k, v)
		}
		return
	}
}

func (kvt KeyValAll[K, V, Out]) FilterOut(f filter.Filter[Out]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		out = kvt(k, v)
		include = f(out)
		return
	}
}

func (vt ValAll[V, Out]) FilterVal(f filter.Filter[V]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		include = f(v)
		if include {
			out = vt(v)
		}
		return
	}
}

func (vt ValAll[V, Out]) FilterOut(f filter.Filter[Out]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		out = vt(v)
		include = f(out)
		return
	}
}

func (kvt KeyVal[K, V, Out]) FilterKey(f filter.Filter[K]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(k)
		if include {
			out, include = kvt(k, v)
		}
		return
	}
}

func (kvt KeyVal[K, V, Out]) FilterVal(f filter.Filter[V]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		include = f(v)
		if include {
			out, include = kvt(k, v)
		}
		return
	}
}

func (kvt KeyVal[K, V, Out]) FilterOut(f filter.Filter[Out]) KeyVal[K, V, Out] {
	return func(k K, v V) (out Out, include bool) {
		out, include = kvt(k, v)
		if include {
			include = f(out)
		}
		return
	}
}

func (vt Val[V, Out]) FilterVal(f filter.Filter[V]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		include = f(v)
		if include {
			out, include = vt(v)
		}
		return
	}
}

func (vt Val[V, Out]) FilterOut(f filter.Filter[Out]) Val[V, Out] {
	return func(v V) (out Out, include bool) {
		out, include = vt(v)
		if include {
			include = f(out)
		}
		return
	}
}
