package type32

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// TODO: this logic should be moved to a bimap
type typeMap struct {
	t2u map[reflect.Type]uint32
	u2t map[uint32]reflect.Type
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

// RegisterType32s registers many TypeIDer32s.
func (tm typeMap) RegisterType32s(zeroValues ...TypeIDer32) {
	for _, zeroValue := range zeroValues {
		tm.Add(reflect.TypeOf(zeroValue), zeroValue.TypeID32())
	}
}
