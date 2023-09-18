package serial

// TypeSerializer takes an interface and returns the serialization as a byte
// slice that includes the type data.
type TypeSerializer interface {
	SerializeType(any, []byte) ([]byte, error)
}

// TypeDeserializer can infer the type from the message. It can take a byte
// slice and return an interface containing an instance to the value represented
// by the serialized data.
type TypeDeserializer interface {
	DeserializeType([]byte) (any, error)
}

// InterfaceDeserializer is what is provided by the json and gob libraries,
// the ability to populate an interface from a byte slice.
type InterfaceDeserializer interface {
	Deserialize(any, []byte) error
}
