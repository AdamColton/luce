package rye

// S provides a namespace to attach the Serialize methods.
type S struct{}

// Serialize holds the Serialize methods.
var Serialize S

// Uint16 encodes x and writes it to b.
func (S) Uint16(b []byte, x uint16) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
}

// Uint32 encodes x and writes it to b.
func (S) Uint32(b []byte, x uint32) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
}

// Uint64 encodes x and writes it to b.
func (S) Uint64(b []byte, x uint64) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
	b[4] = byte(x >> 32)
	b[5] = byte(x >> 40)
	b[6] = byte(x >> 48)
	b[7] = byte(x >> 56)
}
