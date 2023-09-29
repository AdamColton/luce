package type32

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/rye"
)

// MapPrefixer fulfills serial.ReflectTypePrefixer, mapping each type to the
// uint32 that prefixes it.
type MapPrefixer map[reflect.Type]uint32

// PrefixReflectType fulfills ReflectTypePrefixer. It appends the uint32 for the
// type to b as 4 bytes, little-endian. If the type is not in the map, or the map
// is nil, it returns b and an ErrTypeNotFound.
func (p MapPrefixer) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if p == nil {
		return b, ErrTypeNotFound{t}
	}
	u, ok := p[t]
	if !ok {
		return b, ErrTypeNotFound{t}
	}
	ln := len(b)
	b = slice.New(b).CheckCapacity(ln+4, (ln*2)+4)
	b = b[:ln+4]
	rye.Serialize.Uint32(b[ln:ln+4], u)
	return b, nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// MapPrefixer as the InterfaceTypePrefixer and the provided Serializer.
func (p MapPrefixer) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: serial.WrapPrefixer(p),
		Serializer:            s,
	}
}
