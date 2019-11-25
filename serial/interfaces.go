package serial

// Serializer takes an interface and returns the serialization as a byte slice.
type Serializer interface {
	Serialize(any, []byte) ([]byte, error)
}

// Deserializer takes an interface and a serialization of the underlying type
// and populates the interface from the data.
type Deserializer interface {
	Deserialize(any, []byte) error
}
