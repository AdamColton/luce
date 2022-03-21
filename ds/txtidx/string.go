package txtidx

// String is, an indexed string.
type String struct {
	Len uint32
	// If Start points to IndexWords, it points to the index
	// if it points to UnindexWords index = MaxUint32-Start
	Start        uint32
	IndexWords   []IWordIndex
	UnindexWords []UWordIndex
}

// Link the current word to the next word
type Link struct {
	VID  uint16
	Next uint32
}

// UWordIndex
type UWordIndex struct {
	ID   uint32
	Next []uint32
}

// PreParser takes a description of a stirng and
// creates a String from it.
type PreParser struct {
	IndexedCount uint32
	IDXs         []SourceIDX
	*Source
	uniqueRoots    map[uint32]uint32
	sIdxToLocalIdx map[SourceIDX]uint16
	lengths        map[SourceIDX]uint32
}

func (s *Source) NewPreParser() *PreParser {
	return &PreParser{
		Source:         s,
		uniqueRoots:    make(map[uint32]uint32),
		sIdxToLocalIdx: make(map[SourceIDX]uint16),
		lengths:        make(map[SourceIDX]uint32),
	}
}

// Indexed inserts an indexted word
func (p *PreParser) Indexed(root, variant string) {
	idx := p.UpsertIndexed(root, variant)
	r := idx.Root()
	_, found := p.uniqueRoots[r]
	if !found {
		p.IndexedCount++
		p.uniqueRoots[r] = MaxUint32
		p.lengths[idx] = uint32(len([]byte(variant)))
	}
	p.IDXs = append(p.IDXs, idx)
}

// Unindex inserts an unindexed word
func (p *PreParser) Unindex(str string) {
	idx := p.UpsertUnindexed(str)
	p.IDXs = append(p.IDXs, idx)
	p.lengths[idx] = uint32(len([]byte(str)))
}

func (p *PreParser) String() *String {
	s := &String{
		Start: MaxUint32,
	}
	for i := len(p.IDXs) - 1; i >= 0; i-- {
		idx := p.IDXs[i]
		s.Len += p.lengths[idx]
		if idx.Indexed() {
			p.insertIndexed(idx, s)
		} else {
			p.insertUnindexed(idx, s)
		}
	}
	return s
}
func (p *PreParser) insertIndexed(idx SourceIDX, s *String) {
	r := idx.Root()
	iwIdx := p.uniqueRoots[r]
	if iwIdx == MaxUint32 {
		s.IndexWords = append(s.IndexWords, IWordIndex{
			BaseID:   r,
			Variants: []uint16{idx.Variant()},
			Links: []Link{{
				VID:  0,
				Next: s.Start,
			}},
		})
		p.sIdxToLocalIdx[idx] = 0
		s.Start = uint32(len(s.IndexWords)) - 1
		p.uniqueRoots[r] = s.Start
	} else {
		iw := s.IndexWords[iwIdx]
		vid, found := p.sIdxToLocalIdx[idx]
		if !found {
			vid = uint16(len(iw.Variants))
			iw.Variants = append(iw.Variants, idx.Variant())
		}
		iw.Links = append(iw.Links, Link{
			VID:  vid,
			Next: s.Start,
		})
		s.IndexWords[iwIdx] = iw
		s.Start = iwIdx
	}
}

func (p *PreParser) insertUnindexed(idx SourceIDX, s *String) {
	lIdx, found := p.sIdxToLocalIdx[idx]
	if !found {
		s.UnindexWords = append(s.UnindexWords, UWordIndex{
			ID:   idx.Unindex(),
			Next: []uint32{s.Start},
		})
		s.Start = p.IndexedCount + uint32(len(s.UnindexWords)) - 1
	} else {
		uw := s.UnindexWords[lIdx]
		uw.Next = append(uw.Next, s.Start)
		s.Start = p.IndexedCount + uint32(lIdx)
	}
}
