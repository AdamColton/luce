package type32

import (
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

// TypeIDer32 identifies a type by a uint32. The uint32 size was chosen becuase
// it should allow for plenty of TypeID32 types, but uses little overhead.
type TypeIDer32 interface {
	TypeID32() uint32
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
