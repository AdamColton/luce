package serial

import "reflect"

// Serializer appends the serialization of a value to a byte slice and returns
// the result. Passing a nil slice is always safe.
type Serializer interface {
	Serialize(any, []byte) ([]byte, error)
}

// Deserializer takes an interface and a serialization of the underlying type
// and populates the interface from the data.
type Deserializer interface {
	Deserialize(any, []byte) error
}

// InterfaceTypePrefixer appends the type of the value to a byte slice and
// returns the result. Generally this will end up effectively prefixing the
// type.
type InterfaceTypePrefixer interface {
	PrefixInterfaceType(any, []byte) ([]byte, error)
}

// ReflectTypePrefixer appends a reflect.Type to a byte slice and returns the
// result. Generally this will end up effectively prefixing the type.
type ReflectTypePrefixer interface {
	PrefixReflectType(reflect.Type, []byte) ([]byte, error)
}

// TypePrefixer combines both type prefixing techniques.
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
