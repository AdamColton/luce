package serial

import (
	"reflect"
)

// WrapPrefixer takes a ReflectTypePrefixer and wraps it with logic to add
// PrefixInterfaceType, thereby fulfilling TypePrefixer. This makes it easy to
// fulfill any type of prefixer by just fulfilling ReflectTypePrefixer.
func WrapPrefixer(pre ReflectTypePrefixer) TypePrefixer {
	if p, ok := pre.(TypePrefixer); ok {
		return p
	}
	return wrapPrefixer{pre}
}

type wrapPrefixer struct {
	ReflectTypePrefixer
}

// PrefixInterfaceType fulfills InterfaceTypePrefixer. It appends the type of i
// to b, using the underlying ReflectTypePrefixer.
func (wp wrapPrefixer) PrefixInterfaceType(i any, b []byte) ([]byte, error) {
	return wp.PrefixReflectType(reflect.TypeOf(i), b)
}

// Wrapped fulfills upgrade.Wrapper allowing the ReflectTypePrefixer to
// be upgraded to any other types it may fulfill.
func (wp wrapPrefixer) Wrapped() any {
	return wp.ReflectTypePrefixer
}
