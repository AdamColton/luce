package txtidx

// Source holds all the words in the "dictionary"
type Source struct {
	Indexed   []IndexedWord
	Unindexed []string
	Lookup    map[string]SourceIDX
}

// IndexedWord holds all the varaintes of an index word with the word at 0 being
// the root word.
type IndexedWord struct {
	Variants []string
}

func NewSource() *Source {
	return &Source{
		Lookup: make(map[string]SourceIDX),
	}
}

func (s *Source) UpsertRoot(root string) SourceIDX {
	idx, found := s.Lookup[root]
	if found {
		if idx.Variant() != 0 {
			panic("cannot add variant as root")
		}
		return idx
	}
	idx = SourceIDX(len(s.Indexed))
	s.Indexed = append(s.Indexed, IndexedWord{
		Variants: []string{root},
	})
	s.Lookup[root] = idx
	return idx
}

func (s *Source) UpsertIndexed(root, variant string) SourceIDX {
	idx, found := s.Lookup[variant]
	if found {
		return idx
	}
	ridx := s.UpsertRoot(root)
	iw := s.Indexed[ridx]
	idx = ridx.SetVariant(len(iw.Variants))
	iw.Variants = append(iw.Variants, variant)
	s.Indexed[ridx] = iw
	s.Lookup[variant] = idx
	return idx
}

func (s *Source) UpsertUnindexed(str string) SourceIDX {
	idx, found := s.Lookup[str]
	if found {
		return idx
	}
	idx = UnindexedIDX(len(s.Unindexed))
	s.Unindexed = append(s.Unindexed, str)
	s.Lookup[str] = idx
	return idx
}

func (s *Source) Get(idx SourceIDX) string {
	if idx.Indexed() {
		return s.Indexed[idx.Root()].Variants[idx.Variant()]
	}
	return s.Unindexed[idx.Unindex()]
}

// SourceIDX
// 1 bit for Indexed
// 31 bits for variant
// 32 bits for root
type SourceIDX uint64

const (
	isUnindexed SourceIDX = 1 << 63
	rootMask              = (1 << 32) - 1
)

func (s SourceIDX) Indexed() bool {
	return s < isUnindexed
}

func (s SourceIDX) Variant() uint64 {
	return uint64(s) >> 32
}
func (s SourceIDX) SetVariant(v int) SourceIDX {
	return (s & rootMask) | (SourceIDX(v) << 32)
}

func (s SourceIDX) Root() uint64 {
	return uint64(s & rootMask)
}

func (s SourceIDX) Unindex() uint64 {
	return uint64(s &^ isUnindexed)
}

func IndexedIDX(root, variant uint64) SourceIDX {
	return SourceIDX((variant << 32) | root)
}

func UnindexedIDX(idx int) SourceIDX {
	return isUnindexed | SourceIDX(idx)
}
