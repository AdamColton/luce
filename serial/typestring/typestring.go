// Package typestring identifies types by a string so that serialized data can
// be prefixed with its type.
package typestring

import (
	"reflect"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

const (
	// ErrTypeNotFound is returned when a prefixer has no string for a type.
	ErrTypeNotFound = lerr.Str("Type was not found")
	// ErrNilZero is returned when RegisterType is given nil.
	ErrNilZero = lerr.Str("TypeIDStringer.Register) cannot register nil interface")
	// ErrMalformed is returned by GetType when the data does not start with a
	// non-empty type string followed by a space.
	ErrMalformed = lerr.Str("Type string is malformed")
	// ErrNotRegistered is returned by GetType when the type string is not
	// registered.
	ErrNotRegistered = lerr.Str("No type registered")
)

// TypeIDStringer identifies a type by a string. The string must not contain a
// space, because a space ends the type string when it is read.
type TypeIDStringer interface {
	TypeIDString() string
}

// MapPrefixer fulfills serial.ReflectTypePrefixer, mapping each type to the
// string that prefixes it.
type MapPrefixer map[reflect.Type]string

// PrefixReflectType fulfills ReflectTypePrefixer. It appends the string for the
// type followed by a space to b. If the type is not in the map, or the map is
// nil, it returns nil and ErrTypeNotFound.
func (p MapPrefixer) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	if p == nil {
		return nil, ErrTypeNotFound
	}
	s, ok := p[t]
	if !ok {
		return nil, ErrTypeNotFound
	}
	s += " "
	return append(b, []byte(s)...), nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// MapPrefixer as the InterfaceTypePrefixer and the provided Serializer.
func (p MapPrefixer) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: serial.WrapPrefixer(p),
		Serializer:            s,
	}
}

// StringPrefixer fulfills serial.InterfaceTypePrefixer but requires that the
// interfaces passed to it fulfill TypeIDStringer.
type StringPrefixer struct{}

// PrefixInterfaceType casts i to TypeIDStringer and appends its string followed
// by a space to b. If i does not fulfill TypeIDStringer, it returns nil and
// ErrTypeNotFound.
func (StringPrefixer) PrefixInterfaceType(i any, b []byte) ([]byte, error) {
	ts, ok := i.(TypeIDStringer)
	if !ok {
		return nil, ErrTypeNotFound
	}
	s := ts.TypeIDString() + " "
	b = append(b, []byte(s)...)
	return b, nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// StringPrefixer as the InterfaceTypePrefixer and the provided Serializer.
func (t StringPrefixer) Serializer(s serial.Serializer) serial.PrefixSerializer {
	return serial.PrefixSerializer{
		InterfaceTypePrefixer: t,
		Serializer:            s,
	}
}

type typeMap struct {
	t2s map[reflect.Type]string
	s2t map[string]reflect.Type
}

// TypeMap tracks the mapping between types and their string values. It can
// register types, prefix data with a type string and read the type back from
// the data. It has an unexported method so it can only be created with
// NewTypeMap.
type TypeMap interface {
	serial.TypeRegistrar
	serial.TypePrefixer
	serial.Detyper
	Add(t reflect.Type, id string)
	RegisterTypeString(zeroValue TypeIDStringer)
	Serializer(s serial.Serializer) serial.PrefixSerializer
	WriterSerializer(s serial.WriterSerializer) serial.PrefixSerializer
	Deserializer(d serial.Deserializer) serial.PrefixDeserializer
	ReaderDeserializer(d serial.ReaderDeserializer) serial.PrefixDeserializer
	private()
}

// NewTypeMap creates a TypeMap and registers the zeroValues with
// RegisterTypeString. A nil zeroValue panics.
func NewTypeMap(zeroValues ...TypeIDStringer) TypeMap {
	tm := typeMap{
		t2s: make(map[reflect.Type]string),
		s2t: make(map[string]reflect.Type),
	}
	for _, z := range zeroValues {
		tm.RegisterTypeString(z)
	}
	return tm
}

func (typeMap) private() {}

// RegisterType fulfills serial.TypeRegistrar. The zeroValue must fulfill
// TypeIDStringer. It returns ErrNilZero for nil and an error for any value that
// does not fulfill TypeIDStringer.
func (tm typeMap) RegisterType(zeroValue any) error {
	zvStr, ok := zeroValue.(TypeIDStringer)
	if ok {
		tm.RegisterTypeString(zvStr)
		return nil
	}
	if zeroValue == nil {
		return ErrNilZero
	}
	return lerr.Str("TypeIDStringer.Register) " + reflect.TypeOf(zeroValue).Name() + " does not fulfill TypeIDStringer")
}

// RegisterTypeString registers a TypeIDStringer. It functions the same as
// serial.TypeRegistrar but adds type safety. The type of the zeroValue is
// registered, so a pointer registers the pointer type.
func (tm typeMap) RegisterTypeString(zeroValue TypeIDStringer) {
	tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeIDString())
}

// Add maps a type to an id. This allows for types that do not fulfill
// TypeIDStringer to be registered. The id must not contain a space. Adding a
// type or an id that is already registered replaces its mapping.
func (tm typeMap) Add(t reflect.Type, id string) {
	tm.t2s[t] = id
	tm.s2t[id] = t
}

// PrefixReflectType fulfills serial.ReflectTypePrefixer. It appends the id of
// the type followed by a space, or returns ErrTypeNotFound.
func (tm typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	return MapPrefixer(tm.t2s).PrefixReflectType(t, b)
}

// PrefixInterfaceType fulfills serial.InterfaceTypePrefixer. It appends the id
// of the type of i followed by a space, or returns ErrTypeNotFound.
func (tm typeMap) PrefixInterfaceType(i any, b []byte) ([]byte, error) {
	return serial.WrapPrefixer(MapPrefixer(tm.t2s)).PrefixInterfaceType(i, b)
}

// GetType fulfills serial.Detyper. It reads the type string, which ends at the
// first space, and returns the registered type and the data after the space. It
// returns ErrMalformed if there is no space or the type string is empty and
// ErrNotRegistered if the type string is not registered.
func (tm typeMap) GetType(data []byte) (t reflect.Type, rest []byte, err error) {
	s := string(data)
	idx := strings.IndexRune(string(data), ' ')
	if idx < 1 {
		return nil, nil, ErrMalformed
	}
	s = s[:idx]

	rt := tm.s2t[s]
	if rt == nil {
		return nil, nil, ErrNotRegistered
	}

	return rt, data[idx+1:], nil
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
