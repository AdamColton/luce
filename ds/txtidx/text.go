package txtidx

const MaxUint32 = ^uint32(0)

// IWordIndex holds all reference to a root word in the document. By putting all
// the Variants in one place it should be easier to index.
type IWordIndex struct {
	BaseID   uint32
	Variants []uint16
	Links    []Link
}

func Build(s *Source, d *String) string {
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

func NewLookup(s *Source, d *String) *Lookup {
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
