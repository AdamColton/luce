package txtidx

const MaxUint32 = ^uint32(0)

type DocumentIndex struct {
	Len          uint32
	Start        uint32
	IndexWords   []IWordIndex
	UnindexWords []UWordIndex
}

// IWordIndex holds all reference to a root word in the document. By putting all
// the Variants in one place it should be easier to index.
type IWordIndex struct {
	BaseID   uint32
	Variants []uint16
	Links    []Link
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

// Source holds all the words in the "dictionary"
type Source struct {
	Indexed   []IndexedWord
	Unindexed []string
	Lookup    map[string]SourceIDX
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

// IndexedWord holds all the varaintes of an index word with the word at 0 being
// the root word.
type IndexedWord struct {
	Variants []string
}

func Build(s *Source, d *DocumentIndex) string {
	out := make([]byte, 0, d.Len)
	l := NewLookup(s, d)

	cursors := make([]uint16, len(d.IndexWords)+len(d.UnindexWords))
	cur := d.Start
	lniw := uint32(len(d.IndexWords))
	var word string
	for cur != MaxUint32 {
		nc := cursors[cur]
		cursors[cur]++
		if cur >= lniw {
			cur -= lniw
			word = l.Unindexed[cur]
			cur = d.UnindexWords[cur].Next[nc]
		} else {
			i := d.IndexWords[cur]
			lk := i.Links[nc]
			word = l.IndexedVariants[cur][lk.VID]
			cur = lk.Next

		}
		out = append(out, []byte(word)...)
	}
	return string(out)
}

type Lookup struct {
	IndexedVariants [][]string
	Unindexed       []string
}

func NewLookup(s *Source, d *DocumentIndex) *Lookup {
	l := &Lookup{
		IndexedVariants: make([][]string, len(d.IndexWords)),
		Unindexed:       make([]string, len(d.UnindexWords)),
	}

	for i, iwi := range d.IndexWords {
		iw := s.Indexed[iwi.BaseID]
		vs := make([]string, len(iwi.Variants))
		for j, vIdx := range iwi.Variants {
			vs[j] = iw.Variants[vIdx]
		}
		l.IndexedVariants[i] = vs
	}

	for i, u := range d.UnindexWords {
		l.Unindexed[i] = s.Unindexed[u.ID]
	}

	return l
}
