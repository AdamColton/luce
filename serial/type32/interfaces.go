package type32

import (
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

// Sentinel errors
const (
	// ErrTooShort is returned by GetType when the data is shorter than the 4
	// bytes of a type ID.
	ErrTooShort = lerr.Str("TypeID32 too short")
	// ErrNotRegistered is returned by GetType when the type ID is not in the
	// TypeMap.
	ErrNotRegistered = lerr.Str("No type registered")
	// ErrSerNotT32 is the error for serializing a value that does not fulfill
	// TypeIDer32.
	ErrSerNotT32 = lerr.Str("Serialize requires interface to be TypeIDer32")
	// ErrNilZero is returned by RegisterType when it is given nil.
	ErrNilZero = lerr.Str("TypeID32Deserializer.Register) cannot register nil interface")
)

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen because
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
}

// TypeMap tracks the mapping between types and their uint32 values. It can
// register types, prefix data with a type ID (4 bytes, little-endian) and read
// the type back from the data. It is not safe for concurrent use: registering
// writes to its maps without a lock. It has an unexported method, so it can
// only be created with NewTypeMap.
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
