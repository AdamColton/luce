package rye

import "math"

// Deserializer provides a helper for deserializing binary data
type Deserializer struct {
	Data []byte
	Idx  int
}

// NewDeserializer returns a Deserializer prepared to deserialize the provided
// data
func NewDeserializer(data []byte) *Deserializer {
	return &Deserializer{
		Data: data,
	}
}

// Done checks if the Index is at or past the end of the Data.
func (d *Deserializer) Done() bool {
	return d.Idx >= len(d.Data)
}

// Sub returns a Sub-Deserializer of a given length. The index of the parent
// is placed at the end of the data allocated to the Sub-Deserializer.
func (d *Deserializer) Sub(ln int) *Deserializer {
	d.Idx += ln
	return &Deserializer{
		Data: d.Data[d.Idx-ln : d.Idx],
	}
}

// Slice returns the next ln bytes, which share memory with Data, and increases
// the index.
func (d *Deserializer) Slice(ln int) []byte {
	d.Idx += ln
	return d.Data[d.Idx-ln : d.Idx]
}

// Byte returns one byte from the Deserializer and increases the index.
func (d *Deserializer) Byte() byte {
	d.Idx++
	return d.Data[d.Idx-1]
}

// Uint8 returns a uint8 from the Deserializer and increases the index.
func (d *Deserializer) Uint8() uint8 {
	d.Idx++
	return uint8(d.Data[d.Idx-1])
}

// Uint16 returns a uint16 from the Deserializer and increases the index.
func (d *Deserializer) Uint16() uint16 {
	d.Idx += 2
	return Deserialize.Uint16(d.Data[d.Idx-2:])
}

// Uint32 returns a uint32 from the Deserializer and increases the index.
func (d *Deserializer) Uint32() uint32 {
	d.Idx += 4
	return Deserialize.Uint32(d.Data[d.Idx-4:])
}

// Uint64 returns a uint64 from the Deserializer and increases the index.
func (d *Deserializer) Uint64() uint64 {
	d.Idx += 8
	return Deserialize.Uint64(d.Data[d.Idx-8:])
}

// Uint reads a little-endian uint64 of size bytes and increases the index.
func (d *Deserializer) Uint(size byte) uint64 {
	var out uint64
	for i := byte(0); i < size; i++ {
		out |= uint64(d.Byte()) << (i * 8)
	}
	return out
}

// Float32 reads a float32 and increases the index.
func (d *Deserializer) Float32() float32 {
	return math.Float32frombits(d.Uint32())
}

// Float64 reads a float64 and increases the index.
func (d *Deserializer) Float64() float64 {
	return math.Float64frombits(d.Uint64())
}

// Int8 reads an int8 and increases the index.
func (d *Deserializer) Int8() int8 {
	return int8(d.Byte())
}

// Int16 reads an int16 and increases the index.
func (d *Deserializer) Int16() int16 {
	return int16(d.Uint16())
}

// Int32 reads an int32 and increases the index.
func (d *Deserializer) Int32() int32 {
	return int32(d.Uint32())
}

// Int64 reads an int64 and increases the index.
func (d *Deserializer) Int64() int64 {
	return int64(d.Uint64())
}

// String reads a string of ln bytes and increases the index.
func (d *Deserializer) String(ln int) string {
	return string(d.Slice(ln))
}
