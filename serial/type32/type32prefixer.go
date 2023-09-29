package type32

import (
	"reflect"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/rye"
)

// Type32Prefixer fulfills PrefixInterfaceType but requires that the interfaces
// passed to it fulfill TypeIDer32.
type Type32Prefixer struct{}

// PrefixInterfaceType casts i to TypeIDer32 and prefixes 4 bytes with that
// value.
func (Type32Prefixer) PrefixInterfaceType(i any, b []byte) ([]byte, error) {
	t32, ok := i.(TypeIDer32)
	if !ok {
		return nil, ErrTypeNotFound{reflect.TypeOf(i)}
	}
	ln := len(b)
	b = slice.New(b).CheckCapacity(ln+4, (ln*2)+4)
	b = b[:ln+4]
	rye.Serialize.Uint32(b[ln:ln+4], t32.TypeID32())
	return b, nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// Type32Prefixer as the InterfaceTypePrefixer and the provided Serializer.
func (t Type32Prefixer) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: t,
		Serializer:            s,
	}
}
