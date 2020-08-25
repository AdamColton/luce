package rye

// Serializer is used to Serialize into the Data field.
type Serializer struct {
	Data []byte
	Size int
	Idx  int
}

// Make will set Data to the length of Size. If Data is already populated, it
// will be appeded to.
func (s *Serializer) Make() *Serializer {
	ln := s.Size - len(s.Data)
	if ln > 0 {
		s.Data = append(s.Data, make([]byte, ln)...)
	}
	return s
}

// Sub serializer of a specific length.
func (s *Serializer) Sub(ln int) *Serializer {
	s.Idx += ln
	return &Serializer{
		Data: s.Data[s.Idx-ln : s.Idx],
		Size: ln,
	}
}

// Byte writes a byte to the Serializer and increases the index.
func (s *Serializer) Byte(b byte) {
	s.Data[s.Idx] = b
	s.Idx++
}

// Uint8 writes a uint8 to the Serializer and increases the index.
func (s *Serializer) Uint8(x uint8) {
	s.Data[s.Idx] = byte(x)
	s.Idx++
}

// Uint16 writes a uint16 to the Serializer and increases the index.
func (s *Serializer) Uint16(x uint16) {
	s.Data[s.Idx] = byte(x)
	s.Data[s.Idx+1] = byte(x >> 8)
	s.Idx += 2
}

// Uint32 writes a uint32 to the Serializer and increases the index.
func (s *Serializer) Uint32(x uint32) {
	s.Data[s.Idx] = byte(x)
	s.Data[s.Idx+1] = byte(x >> 8)
	s.Data[s.Idx+2] = byte(x >> 16)
	s.Data[s.Idx+3] = byte(x >> 24)
	s.Idx += 4
}

// Uint64 writes a uint64 to the Serializer and increases the index.
func (s *Serializer) Uint64(x uint64) {
	s.Data[s.Idx] = byte(x)
	s.Data[s.Idx+1] = byte(x >> 8)
	s.Data[s.Idx+2] = byte(x >> 16)
	s.Data[s.Idx+3] = byte(x >> 24)
	s.Data[s.Idx+4] = byte(x >> 32)
	s.Data[s.Idx+5] = byte(x >> 40)
	s.Data[s.Idx+6] = byte(x >> 48)
	s.Data[s.Idx+7] = byte(x >> 56)
	s.Idx += 8
}

// Uint serializes the value little-endian. If size is >0, that number of bytes
// will be written. If size==0, it will write the value with fewest bytes,
// omitting leading zeros. The returned byte indicates the number of bytes
// written. Size should be between 0 and 8, but it is not checked. This does
// not increase the Idx.
func (s *Serializer) Uint(size byte, value uint64) byte {
	if size == 0 && value == 0 {
		s.Byte(0)
		return 1
	}
	var out byte
	for (size > 0 && out < size) || (size == 0 && value > 0) {
		s.Byte(byte(value))
		value >>= 8
		out++
	}
	return out
}

// Int8 writes a int8 to the Serializer and increases the index.
func (s *Serializer) Int8(x int8) {
	s.Uint8(uint8(x))
}

// Int16 writes a int16 to the Serializer and increases the index.
func (s *Serializer) Int16(x int16) {
	s.Uint16(uint16(x))
}

// Int32 writes a int32 to the Serializer and increases the index.
func (s *Serializer) Int32(x int32) {
	s.Uint32(uint32(x))
}

// Int64 writes a int64 to the Serializer and increases the index.
func (s *Serializer) Int64(x int64) {
	s.Uint64(uint64(x))
}

// Float32 writes a float32 to the Serializer and increases the index.
func (s *Serializer) Float32(f float32) {
	s.Uint32(float32ToUint32(f))
}

// Float64 writes a float64 to the Serializer and increases the index.
func (s *Serializer) Float64(f float64) {
	s.Uint64(float64ToUint64(f))
}

// Slice writes a byte slice to the Serializer and increases the index.
func (s *Serializer) Slice(data []byte) {
	copy(s.Data[s.Idx:], data)
	s.Idx += len(data)
}

// String writes a string to the Serializer and increases the index.
func (s *Serializer) String(data string) {
	s.Slice([]byte(data))
}
