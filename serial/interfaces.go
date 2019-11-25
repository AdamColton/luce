package serial

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
