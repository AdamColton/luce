package serial

import "reflect"

// Serializer takes an interface and returns the serialization as a byte slice.
type Serializer interface {
	Serialize(interface{}, []byte) ([]byte, error)
}

// Deserializer takes an interface and a serialization of the underlying type
// and populates the interface from the data.
type Deserializer interface {
	Deserialize(interface{}, []byte) error
}

// InterfaceTypePrefixer writes the type of the interface to slice. Generally
// this will end up effectivly prefixing the type.
type InterfaceTypePrefixer interface {
	PrefixInterfaceType(interface{}, []byte) ([]byte, error)
}

// ReflectTypePrefixer writes a reflect.Type to slice. Generally this will end up
// effectivly prefixing the type.
type ReflectTypePrefixer interface {
	PrefixReflectType(reflect.Type, []byte) ([]byte, error)
}

// TypePrefixer combines both TypePrefixing techniques.
type TypePrefixer interface {
	ReflectTypePrefixer
	InterfaceTypePrefixer
}
