package byteid

// Index a slice by a byte ID. Allows the equivalent of map[[]byte]<Type>. A
// slice of the desired type is kept seperatly and the Index values are managed.
type Index interface {
	// Insert an ID. The first value returned is the index and the bool
	// indicates if an append is required.
	Insert(id []byte) (int, bool)
	// Get by ID. If not found it should return (-1,false). If it is found the
	// first value is the index and the second value is True.
	Get(id []byte) (int, bool)
	// Delete by ID. Removes the ID from the index, the value should be
	// recycled. This should be called before removing the value from the slice.
	Delete(id []byte) (int, bool)
	// SliceLen of the Indexed slice.
	SliceLen() int
	// SetSliceLen can be used to grow the slice.
	SetSliceLen(int)
	// Next ID after the ID given
	Next(id []byte) ([]byte, int)
}

type IndexFactory func(slicelen int) Index
