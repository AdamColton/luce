// Package upgrade solves an issue that results from the intersections of two
// common patterns in luce. There are many interface wrappers (or decorators).
// These wrap an interface to add functionality. The other pattern is
// upgradeable interfaces. This is were an object can provide additional
// functionality and hinting by fulfilling optional interfaces.
//
// For example, take the following
//   - bytes.Buffer: fulfills io.Writer and io.StringWriter
//   - WriterWrapper: wraps an io.Writer, does not fulfill io.StringWriter
//   - func Foo(w io.Writer): tries to cast w to StringWriter, uses fallback code if it can't
//
// If a Buffer is passed into Foo, the StringWriter cast will work, but if a
// WriterWrapper around a Buffer is passed in, it will fail.
package upgrade

// Wrapper allows an object fulling the Wrapper pattern to return the underlying
// wrapped object. This allows the To function to upgade the underlying object.
type Wrapper interface {
	Wrapped() any
}

// To checks if an object fulfills T. If it does not and the object fulfills
// Wrapper, it also checks if the wrapped object fulfills T.
func To[T any](wrapper any) (to T, ok bool) {
	i := wrapper
	for {
		st, stOk := i.(T)
		if stOk {
			ok = true
			to = st
		}
		sw, swOk := i.(Wrapper)
		if !swOk {
			break
		}
		i = sw.Wrapped()
	}
	return
}
