package txtidx

type rIdx struct {
	idx   int
	vIdxs map[VarID]int
}

type preParser struct {
	*Corpus
	start   []byte
	words   []string
	rootIdx map[WordID]*rIdx

	*Document
}

func (c *Corpus) newPP() *preParser {
	return &preParser{
		Corpus:  c,
		rootIdx: map[WordID]*rIdx{},
	}
}

func (pp *preParser) build() *Document {
	pp.Document = &Document{
		DocID: pp.Max.DocID,
		start: pp.start,
		Ln:    uint32(len(pp.start)),
	}
	pp.Max.DocID++

	for _, w := range pp.words {
		pp.buildWord(w)
	}
	return pp.Document
}

func (pp *preParser) buildWord(word string) {
	rootID, vid := pp.Upsert(word)
	rootIdx, found := pp.rootIdx[rootID]
	if !found {
		pp.Corpus.IDs[rootID].Documents.add(pp.DocID)
		rootIdx = &rIdx{
			idx:   len(pp.Words),
			vIdxs: map[VarID]int{},
		}
		pp.rootIdx[rootID] = rootIdx
		pp.Words = append(pp.Words, DocWord{
			WordID: rootID,
		})
	}
	dw := &(pp.Words[rootIdx.idx])
	idx, found := rootIdx.vIdxs[vid]
	if !found {
		idx = len(dw.Variants)
		rootIdx.vIdxs[vid] = idx
		dw.Variants = append(dw.Variants, DocVar{
			VarID: vid,
		})
	}
	dv := &(dw.Variants[idx])
	dv.Positions = append(dv.Positions, pp.Ln)
	pp.Ln += uint32(len(word))
}

func (pp *preParser) set(str string) {
	s := newScanner(str)
	s.matchLetterNumber(true)
	pp.start = s.s[:s.idx]

	for !s.done() {
		start := s.idx
		s.matchLetterNumber(false)
		s.matchLetterNumber(true)
		pp.words = append(pp.words, s.str(start, s.idx))
	}
}
