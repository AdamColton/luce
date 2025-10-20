package reflector

import (
	"reflect"
	"regexp"
)

type TypeMap map[string]reflect.Type

type TypeCollection map[reflect.Type]TypeMap

func TMAdd[T any](key string, tm TypeMap) {
	tm[key] = Type[T]()
}

var embedRe = regexp.MustCompile(`^\w+`)

func TMEmbed[T any](tm TypeMap) {
	t := Type[T]()
	n := t.Name()
	n = embedRe.FindString(n)
	tm[n] = t
}

func TMGet[T any](tmc TypeCollection) TypeMap {
	return tmc[Type[T]()]
}
