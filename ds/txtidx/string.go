package txtidx

type String struct {
	Len          uint32
	Start        uint32
	IndexWords   []IWordIndex
	UnindexWords []UWordIndex
}

// Link the current word to the next word
type Link struct {
	VID  uint8
	Next uint32
}

// UWordIndex
type UWordIndex struct {
	ID   uint32
	Next []uint32
}
