// Package upgrade solves an issue that results from the intersections of two
// common patterns in luce. There are many interface wrappers. These wrap an
// interface to add functionality. The other pattern is upgradeable interfaces.
// This is were an object can provide additional functionality and hinting by
// fulfilling optional interfaces.
//
// However, if an interface that is upgradeable to optional interfaces is
// wrapped, the wrapper will not fulfill those optional interfaces. The Upgrader
// interface is an optional interface that allows wrappers to expose their
// underlying interfaces for upgrade.
package upgrade

import "reflect"

// Upgrader can be implemented by Wrappers to allow the wrapped interface to be
// upgraded.
type Upgrader interface {
	Upgrade(t reflect.Type) interface{}
}

// Wrapped should be called by the Upgrade method by any wrapper fulfilling
// Upgrader.
func Wrapped(i any, t reflect.Type) interface{} {
	if reflect.ValueOf(i).Type().Implements(t) {
		return i
	}
	if u, ok := i.(Upgrader); ok {
		return u.Upgrade(t)
	}
	return nil
}

// Upgrade checks if i fullfils the upgrade type T. T should be an interface.
func Upgrade[T any](i any, t *T) bool {
	v := reflect.ValueOf(t).Type().Elem()
	out := Wrapped(i, v)
	if out == nil {
		return false
	}
	*t = out.(T)
	return true
}
