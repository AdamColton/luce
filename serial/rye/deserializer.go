package rye

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

// Sub returns a Sub-Deserializer of a given length. The index of the parent
// is placed at the end of the data allocated to the Sub-Deserializer.
func (d *Deserializer) Sub(ln int) *Deserializer {
	d.Idx += ln
	return &Deserializer{
		Data: d.Data[d.Idx-ln : d.Idx],
	}
}

// Byte returns one byte from the Deserializer and increases the index.
func (d *Deserializer) Byte() byte {
	d.Idx++
	return d.Data[d.Idx-1]
}

// Uint8 returns a uint8 from the Deserializer and increases the index.
func (d *Deserializer) Uint8() uint8 {
	d.Idx += 1
	return uint8(d.Data[d.Idx-1])
}

// Uint16 returns a uint16 from the Deserializer and increases the index.
func (d *Deserializer) Uint16() uint16 {
	d.Idx += 2
	return uint16(d.Data[d.Idx-1])<<8 + uint16(d.Data[d.Idx-2])
}

func (d *Deserializer) Slice(ln int) []byte {
	d.Idx += ln
	return d.Data[d.Idx-ln : d.Idx]
}

// Uint32 returns a uint32 from the Deserializer and increases the index.
func (d *Deserializer) Uint32() uint32 {
	d.Idx += 4
	return uint32(d.Data[d.Idx-1])<<24 +
		uint32(d.Data[d.Idx-2])<<16 +
		uint32(d.Data[d.Idx-3])<<8 +
		uint32(d.Data[d.Idx-4])
}

// Uint64 returns a uint64 from the Deserializer and increases the index.
func (d *Deserializer) Uint64() uint64 {
	d.Idx += 8
	return uint64(d.Data[d.Idx-1])<<56 +
		uint64(d.Data[d.Idx-2])<<48 +
		uint64(d.Data[d.Idx-3])<<40 +
		uint64(d.Data[d.Idx-4])<<32 +
		uint64(d.Data[d.Idx-5])<<24 +
		uint64(d.Data[d.Idx-6])<<16 +
		uint64(d.Data[d.Idx-7])<<8 +
		uint64(d.Data[d.Idx-8])
}

func (d *Deserializer) Uint(size byte) uint64 {
	var out uint64
	for i := byte(0); i < size; i++ {
		out |= uint64(d.Byte()) << (i * 8)
	}
	return out
}

func (d *Deserializer) Float32() float32 {
	return uint32ToFloat32(d.Uint32())
}

func (d *Deserializer) Float64() float64 {
	return uint64ToFloat64(d.Uint64())
}

func (d *Deserializer) Int8() int8 {
	return int8(d.Byte())
}

func (d *Deserializer) Int16() int16 {
	return int16(d.Uint16())
}

func (d *Deserializer) Int32() int32 {
	return int32(d.Uint32())
}

func (d *Deserializer) Int64() int64 {
	return int64(d.Uint64())
}

func (d *Deserializer) String(ln int) string {
	return string(d.Slice(ln))
}
