package txtidx

const MaxUint32 = ^uint32(0)

type DocumentIndex struct {
	Len          uint32
	Start        uint32
	IndexWords   []IWordIndex
	UnindexWords []UWordIndex
}

type IWordIndex struct {
	BaseID   uint32
	Variants []uint16
	Links    []Link
}

// Link to the next word
type Link struct {
	VID  uint8
	Next uint32
}

type UWordIndex struct {
	ID   uint32
	Next []uint32
}

// Source holds all the words in the "dictionary"
type Source struct {
	Indexed        []IndexedWord
	IndexeLookup   map[string]uint32
	Unindexed      []string
	UnindexeLookup map[string]uint32
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
