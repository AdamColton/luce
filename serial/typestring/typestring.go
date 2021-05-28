package typestring

import (
	"reflect"
	"strings"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

const (
	ErrTypeNotFound  = lerr.Str("Type was not found")
	ErrNilZero       = lerr.Str("TypeIDStringer.Register) cannot register nil interface")
	ErrMalformed     = lerr.Str("Type string is malformed")
	ErrNotRegistered = lerr.Str("No type registered")
)

// TypeIDString identifies a type by a string.
type TypeIDStringer interface {
	TypeIDString() string
}

// MapPrefixer fulfills serial.ReflectTypePrefixer.
type MapPrefixer map[reflect.Type]string

// PrefixReflectType fulfills ReflectTypePrefixer. It will prefix with a word
// followed by a space.
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

// StringPrefixer fulfills PrefixInterfaceType but requires that the interfaces
// passed to it fulfill TypeIDStringer.
type StringPrefixer struct{}

// PrefixInterfaceType casts i to TypeIDer32 and prefixes 4 bytes with that
// value.
func (StringPrefixer) PrefixInterfaceType(i interface{}, b []byte) ([]byte, error) {
	ts, ok := i.(TypeIDStringer)
	if !ok {
		return nil, ErrTypeNotFound
	}
	s := ts.TypeIDString() + " "
	b = append(b, []byte(s)...)
	return b, nil
}

// Serializer is a helper that will create serial.PrefixSerializer using
// Type32Prefixer as the InterfaceTypePrefixer and the provided Serializer.
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

// TypeMap tracks the mapping between types and their uint32 values.
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

// NewTypeMap creates a TypeMap.
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
// TypeIDStringer.
func (tm typeMap) RegisterType(zeroValue interface{}) error {
	zv32, ok := zeroValue.(TypeIDStringer)
	if ok {
		tm.RegisterTypeString(zv32)
		return nil
	}
	if zeroValue == nil {
		return ErrNilZero
	}
	return lerr.Str("TypeIDStringer.Register) " + reflect.TypeOf(zeroValue).Name() + " does not fulfill TypeID32Type")
}

// RegisterTypeString registers a TypeIDStringer. It functions the same as
// serial.TypeRegistrar but adds type safety.
func (tm typeMap) RegisterTypeString(zeroValue TypeIDStringer) {
	tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeIDString())
}

// Add maps a type to an id. This allows for types that do not fulfill
// TypeIDer32 to be registered.
func (tm typeMap) Add(t reflect.Type, id string) {
	tm.t2s[t] = id
	tm.s2t[id] = t
}

// PrefixReflectType fulfills serial.ReflectTypePrefixer.
func (tm typeMap) PrefixReflectType(t reflect.Type, b []byte) ([]byte, error) {
	return MapPrefixer(tm.t2s).PrefixReflectType(t, b)
}

// PrefixInterfaceType fulfills serial.InterfaceTypePrefixer.
func (tm typeMap) PrefixInterfaceType(i interface{}, b []byte) ([]byte, error) {
	return serial.WrapPrefixer(MapPrefixer(tm.t2s)).PrefixInterfaceType(i, b)
}

// GetType fulfills serial.Detyper.
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
// serial.WriterSerializer fulfills serial.Deserializer.
func (tm typeMap) ReaderDeserializer(d serial.ReaderDeserializer) serial.PrefixDeserializer {
	return tm.Deserializer(d)
}
