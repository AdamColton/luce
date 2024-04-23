package entity

// Entity to be written to Store by EntKey.
type Entity interface {
	EntKey() []byte
}

// EntKeySetter can optionally be fulfilled by an Entity. If it is fulfilled, it
// will be invoked whenever an Entity is deserialized from a Store.
type EntKeySetter interface {
	SetEntKey([]byte)
}
