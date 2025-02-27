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

// Detyper takes in serialized data and returns the type of the data and the
// rest of the data (minus the type information).
type Detyper interface {
	GetType(data []byte) (t reflect.Type, rest []byte, err error)
}

// TypeRegistrar is generally required for automatic deserialization. A
// zeroValue is provided (for instance a nil pointer) to register a type that
// can then be deserialized.
type TypeRegistrar interface {
	RegisterType(zeroValue interface{}) error
}

func RegisterPtr[T any](tr TypeRegistrar) error {
	return tr.RegisterType((*T)(nil))
}

// RegisterTypes is a helper to register multiple types in one call.
func RegisterTypes(typeRegistrar TypeRegistrar, zeroValues ...interface{}) error {
	for _, z := range zeroValues {
		err := typeRegistrar.RegisterType(z)
		if err != nil {
			return err
		}
	}
	return nil
}
