package type32

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// Sentinal Errors
const (
	ErrTooShort      = lerr.Str("TypeID32 too short")
	ErrNotRegistered = lerr.Str("No type registered")
	ErrSerNotT32     = lerr.Str("Serialize requires interface to be TypeIDer32")
	ErrNilZero       = lerr.Str("TypeID32Deserializer.Register) cannot register nil interface")
)

type ErrTypeNotFound struct {
	reflect.Type
}

func (err ErrTypeNotFound) Error() string {
	return fmt.Sprintf("Type %s was not found", err.Type)
}

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen becuase
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
}

func uint32ToSlice(u uint32) []byte {
	return []byte{
		byte(u),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 24),
	}
}

func sliceToUint32(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}
	return uint32(b[0]) + (uint32(b[1]) << 8) + (uint32(b[2]) << 16) + (uint32(b[3]) << 24)
}

// MapPrefixer fulfills serial.ReflectTypePrefixer.
type MapPrefixer map[reflect.Type]uint32

// PrefixReflectType fulfills ReflectTypePrefixer. It will prefix with 4 bytes.
func (p MapPrefixer) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if p == nil {
		return nil, ErrTypeNotFound{t}
	}
	u, ok := p[t]
	if !ok {
		return nil, ErrTypeNotFound{t}
	}
	return append(b, uint32ToSlice(u)...), nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// MapPrefixer as the InterfaceTypePrefixer and the provided Serializer.
func (p MapPrefixer) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: serial.WrapPrefixer(p),
		Serializer:            s,
	}
}

// Type32Prefixer fulfills PrefixInterfaceType but requires that the interfaces
// passed to it fulfill TypeIDer32.
type Type32Prefixer struct{}

// PrefixInterfaceType casts i to TypeIDer32 and prefixes 4 bytes with that
// value.
func (Type32Prefixer) PrefixInterfaceType(i interface{}, b []byte) ([]byte, error) {
	t32, ok := i.(TypeIDer32)
	if !ok {
		return nil, ErrTypeNotFound{reflect.TypeOf(i)}
	}
	b = append(b, uint32ToSlice(t32.TypeID32())...)
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

type typeMap struct {
	t2u map[reflect.Type]uint32
	u2t map[uint32]reflect.Type
}

// TypeMap tracks the mapping between types and their uint32 values.
type TypeMap interface {
	serial.TypeRegistrar
	serial.TypePrefixer
	serial.Detyper
	Add(t reflect.Type, id uint32)
	RegisterType32(zeroValue TypeIDer32)
	RegisterType32s(zeroValues ...TypeIDer32)
	Serializer(s serial.Serializer) serial.PrefixSerializer
	WriterSerializer(s serial.WriterSerializer) serial.PrefixSerializer
	Deserializer(d serial.Deserializer) serial.PrefixDeserializer
	ReaderDeserializer(d serial.ReaderDeserializer) serial.PrefixDeserializer
	private()
}

// NewTypeMap creates a TypeMap.
func NewTypeMap() TypeMap {
	return typeMap{
		t2u: make(map[reflect.Type]uint32),
		u2t: make(map[uint32]reflect.Type),
	}
}

func (typeMap) private() {}

// RegisterType fulfills serial.TypeRegistrar. The zeroValue must fulfill
// TypeIDer32.
func (tm typeMap) RegisterType(zeroValue interface{}) error {
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
// serial.TypeRegistrar but adds type safety.
func (tm typeMap) RegisterType32(zeroValue TypeIDer32) {
	tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeID32())
}

// RegisterType32s registers many TypeIDer32s.
func (tm typeMap) RegisterType32s(zeroValues ...TypeIDer32) {
	for _, zeroValue := range zeroValues {
		tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeID32())
	}
}

// Add maps a type to an id. This allows for types that do not fulfill
// TypeIDer32 to be registered.
func (tm typeMap) Add(t reflect.Type, id uint32) {
	tm.t2u[t] = id
	tm.u2t[id] = t
}

// PrefixReflectType fulfills serial.ReflectTypePrefixer.
func (tm typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	return MapPrefixer(tm.t2u).PrefixReflectType(t, b)
}

// PrefixInterfaceType fulfills serial.InterfaceTypePrefixer.
func (tm typeMap) PrefixInterfaceType(i interface{}, b []byte) ([]byte, error) {
	return serial.WrapPrefixer(MapPrefixer(tm.t2u)).PrefixInterfaceType(i, b)
}

// GetType fulfills serial.Detyper.
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

func (tm typeMap) ReaderDeserializer(d serial.ReaderDeserializer) serial.PrefixDeserializer {
	return tm.Deserializer(d)
}
