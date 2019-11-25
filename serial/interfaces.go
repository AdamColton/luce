package serial

// Serializer takes an interface and returns the serialization as a byte slice.
type Serializer interface {
	Serialize(i interface{}) ([]byte, error)
}

// Deserializer knows how to deserialize all types registered with the
// RegisterType method. It can take a byte slice and return an interface
// containing an instance to the value represented by the serialized data.
type Deserializer interface {
	Deserialize([]byte) (interface{}, error)
	RegisterType(zeroValue interface{}) error
}

type ByteID interface {
	Type() []byte
	ID() []byte
}
