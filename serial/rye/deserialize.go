package rye

// D provides a namespace to attach the Deserialize methods.
type D struct{}

// Deserialize holds the Deserialize methods.
var Deserialize D

// Uint16 decodes b as a uint16
func (D) Uint16(b []byte) uint16 {
	return uint16(b[1])<<8 + uint16(b[0])
}

// Uint32 decodes b as a uint32
func (D) Uint32(b []byte) uint32 {
	return uint32(b[3])<<24 +
		uint32(b[2])<<16 +
		uint32(b[1])<<8 +
		uint32(b[0])
}

// Uint64 decodes b as a uint64
func (D) Uint64(b []byte) uint64 {
	return uint64(b[7])<<56 +
		uint64(b[6])<<48 +
		uint64(b[5])<<40 +
		uint64(b[4])<<32 +
		uint64(b[3])<<24 +
		uint64(b[2])<<16 +
		uint64(b[1])<<8 +
		uint64(b[0])
}
