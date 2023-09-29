package type32

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// TODO: this logic should be moved to a bimap
// but I need access to the underlying map for MapPrefixer.
// So I should either expose the maps in bimap
// or in lmap create a map reader and have bimap
// return instances that fulfill that.
type typeMap struct {
	t2u map[reflect.Type]uint32
	u2t map[uint32]reflect.Type
}

// NewTypeMap creates an empty TypeMap. It is not safe for concurrent use.
func NewTypeMap() TypeMap {
	return typeMap{
		t2u: make(map[reflect.Type]uint32),
		u2t: make(map[uint32]reflect.Type),
	}
}

func (typeMap) private() {}

// RegisterType fulfills serial.TypeRegistrar. The zeroValue must fulfill
// TypeIDer32. It returns ErrNilZero for nil and an error for any value that does
// not fulfill TypeIDer32.
func (tm typeMap) RegisterType(zeroValue any) error {
	zv32, ok := zeroValue.(TypeIDer32)
	if ok {
		tm.RegisterType32(zv32)
		return nil
	}
	if zeroValue == nil {
		return ErrNilZero
	}
	return lerr.Str("TypeID32Deserializer.Register) " + reflect.TypeOf(zeroValue).Name() + " does not fulfill TypeID32Type")
}

// RegisterType32 registers a TypeIDer32. It functions the same as
// serial.TypeRegistrar but adds type safety. The type of the zeroValue is
// registered, so a pointer registers the pointer type.
func (tm typeMap) RegisterType32(zeroValue TypeIDer32) {
	tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeID32())
}

// Add maps a type to an id. This allows for types that do not fulfill
// TypeIDer32 to be registered. Adding a type or an id that is already registered
// replaces its mapping.
func (tm typeMap) Add(t reflect.Type, id uint32) {
	tm.t2u[t] = id
	tm.u2t[id] = t
}

// PrefixReflectType fulfills serial.ReflectTypePrefixer. It appends the id of
// the type as 4 bytes, little-endian, or returns an ErrTypeNotFound.
func (tm typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	return MapPrefixer(tm.t2u).PrefixReflectType(t, b)
}

// PrefixInterfaceType fulfills serial.InterfaceTypePrefixer. It appends the id
// of the type of i as 4 bytes, little-endian, or returns an ErrTypeNotFound.
func (tm typeMap) PrefixInterfaceType(i any, b []byte) ([]byte, error) {
	return serial.WrapPrefixer(MapPrefixer(tm.t2u)).PrefixInterfaceType(i, b)
}

// GetType fulfills serial.Detyper. It reads the type ID from the first 4 bytes,
// little-endian, and returns the registered type and the rest of the data. It
// returns ErrTooShort if there are fewer than 4 bytes and ErrNotRegistered if
// the ID is not registered.
func (tm typeMap) GetType(data []byte) (t reflect.Type, rest []byte, err error) {
	if len(data) < 4 {
		return nil, nil, ErrTooShort
	}

	rt := tm.u2t[sliceToUint32(data)]
	if rt == nil {
		return nil, nil, ErrNotRegistered
	}

	return rt, data[4:], nil
}

// Serializer is a helper that will create serial.PrefixSerializer using TypeMap
// as the InterfaceTypePrefixer and the provided Serializer.
func (tm typeMap) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: tm,
		Serializer:            s,
	}
}

// WriterSerializer accepts a WriterSerializer func, which automatically casts
// it to that type so it can be passed into Serializer because
// serial.WriterSerializer fulfills serial.Serializer.
func (tm typeMap) WriterSerializer(s serial.WriterSerializer) serial.PrefixSerializer {
	return tm.Serializer(s)
}

// Deserializer is a helper that will create serial.PrefixDeserializer using
// TypeMap as the Detyper and the provided Deserializer.
func (tm typeMap) Deserializer(d serial.Deserializer) serial.PrefixDeserializer {
	return serial.PrefixDeserializer{
		Detyper:      tm,
		Deserializer: d,
	}
}

// ReaderDeserializer accepts a ReaderDeserializer func, which automatically
// casts it to that type so it can be passed into Deserializer because
// serial.ReaderDeserializer fulfills serial.Deserializer.
func (tm typeMap) ReaderDeserializer(d serial.ReaderDeserializer) serial.PrefixDeserializer {
	return tm.Deserializer(d)
}

// RegisterType32s registers many TypeIDer32s.
func (tm typeMap) RegisterType32s(zeroValues ...TypeIDer32) {
	for _, zeroValue := range zeroValues {
		tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeID32())
	}
}
