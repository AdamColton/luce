package reflector

import (
	"reflect"
	"regexp"
)

// TypeMap maps names to types, for example the names of the fields of a struct
// to their types.
type TypeMap map[string]reflect.Type

// TypeCollection maps a type to its TypeMap.
type TypeCollection map[reflect.Type]TypeMap

// TMAdd adds the type T to tm under key.
func TMAdd[T any](key string, tm TypeMap) {
	tm[key] = Type[T]()
}

var embedRe = regexp.MustCompile(`^\w+`)

// TMEmbed adds the type T to tm under the first word of its name, so an
// embedded field is found under the name of its type without the type
// parameters: slice.Slice[string] is added as "Slice".
func TMEmbed[T any](tm TypeMap) {
	t := Type[T]()
	n := t.Name()
	n = embedRe.FindString(n)
	tm[n] = t
}

// TMGet returns the TypeMap of the type T in tmc, or nil if it has none.
func TMGet[T any](tmc TypeCollection) TypeMap {
	return tmc[Type[T]()]
}
