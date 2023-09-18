package serial

// PrefixSerializer combines a Prefixer and a Serializer to create a
// TypeSerializer.
type PrefixSerializer struct {
	InterfaceTypePrefixer
	Serializer
}

// SerializeType fulfills TypeSerializer using the underlying
// InterfaceTypePrefixer to prefix a type and the Serializer to serialize the
// data.
func (ps PrefixSerializer) SerializeType(i interface{}, b []byte) ([]byte, error) {
	var err error
	b, err = ps.PrefixInterfaceType(i, b)
	if err != nil {
		return nil, err
	}

	return ps.Serialize(i, b)
}
