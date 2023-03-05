// Package upgrade solves an issue that results from the intersections of two
// common patterns in luce. There are many interface wrappers (or decorators).
// These wrap an interface to add functionality. The other pattern is
// upgradeable interfaces. This is where an object can provide additional
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

// Wrapper is fulfilled by an object that wraps another object and can return it.
// This allows the To function to upgrade to an interface that the wrapped
// object fulfills.
type Wrapper interface {
	Wrapped() any
}

// To checks if an object, or any object that it wraps, fulfills T. It follows
// Wrapper.Wrapped through the whole chain of wrappers. If more than one level
// fulfills T, the one that is wrapped the deepest is returned. The bool is
// false if no level fulfills T. A chain of wrappers that loops back on itself
// makes To loop forever.
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
