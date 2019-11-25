// Package idx defines Index, which gives a slice the lookup of a map by mapping
// keys to positions in a slice that the caller keeps separately.
package idx

// Index maps keys to positions in a slice, allowing the equivalent of
// map[Key]<Type>. A slice of the desired type is kept separately by the caller
// and the Index manages which position each key uses, recycling the positions
// of deleted keys.
type Index[Key any] interface {
	// Insert an ID. The first value returned is the index and the bool
	// indicates if an append is required, meaning the index is past the end of
	// the slice. If the ID is already present, its existing index is returned
	// with false.
	Insert(id Key) (int, bool)
	// Get by ID. If not found it should return (-1,false). If it is found the
	// first value is the index and the second value is true.
	Get(id Key) (int, bool)
	// Delete by ID. Removes the ID from the index and returns its index and
	// true, or (-1,false) if the ID was not found. The index is recycled and
	// will be handed out by a later Insert. This should be called before
	// removing the value from the slice.
	Delete(id Key) (int, bool)
	// SliceLen of the Indexed slice: the length it needs to hold every index
	// handed out so far.
	SliceLen() int
	// SetSliceLen raises SliceLen to the given length, for use when the slice
	// was grown or already has values. It never lowers SliceLen.
	SetSliceLen(int)
	// Next returns the smallest ID greater than the ID given, and its index.
	// The ID given does not have to be in the index. If there is no greater
	// ID it returns the zero Key and -1.
	Next(id Key) (Key, int)
	// Len fulfills slice.Lener. Indicates how many values are currently stored
	// in the index.
	Len(int)
}

// IndexFactory creates an empty Index for a slice that has a length of
// slicelen.
type IndexFactory[Key any] func(slicelen int) Index[Key]
