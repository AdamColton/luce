package type32

import (
	"reflect"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/rye"
)

// MapPrefixer fulfills serial.ReflectTypePrefixer.
type MapPrefixer map[reflect.Type]uint32

// PrefixReflectType fulfills ReflectTypePrefixer. It will prefix with 4 bytes.
func (p MapPrefixer) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if p == nil {
		return b, ErrTypeNotFound{t}
	}
	u, ok := p[t]
	if !ok {
		return b, ErrTypeNotFound{t}
	}
	ln := len(b)
	b = checkLn(b)
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
