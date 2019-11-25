package json32

import "github.com/adamcolton/luce/serial/type32"

// Serializer returns a type32.SerializeTypeID32Func that uses the type32 logic
// for prefixing and Serialize for encoding
func Serializer() type32.SerializeTypeID32Func {
	return type32.SerializeTypeID32Func(Serialize)
}

// Deserializer returns a type32.TypeID32Deserializer that uses the type32 logic
// for prefixing and Deserialize for encoding
func Deserializer() *type32.TypeID32Deserializer {
	return type32.DeserializeTypeID32Func(Deserialize).NewTypeID32Deserializer()
}
