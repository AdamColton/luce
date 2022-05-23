package txtidx

import (
	"strings"
)

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	Roots map[string]*Word
	IDs   map[WordID]*Word
	Max   struct {
		WordID
		DocID
	}
}

func NewCorpus() *Corpus {
	return &Corpus{
		Roots: map[string]*Word{},
		IDs:   map[WordID]*Word{},
	}
}

type sig struct{}

type Word struct {
	WordID
	Variants  [][]byte
	VByIDX    map[string]VIDX
	Documents map[DocID]sig
}

type WordID uint32
type VIDX uint32

type DocRef struct {
	VarRefs []VarRef
}

type VarRef struct {
	VIDX
	Locs []uint32
}

func (c *Corpus) Get(word string) (WordID, VIDX) {
	wid, vidx := WordID(MaxUint32), VIDX(MaxUint32)
	rt := root(word)
	w, found := c.Roots[rt]
	if found {
		wid = w.WordID
		v, found := w.VByIDX[word]
		if found {
			vidx = v
		}
	}

	return wid, vidx
}

func (c *Corpus) Upsert(word string) (WordID, VIDX) {
	rt := root(word)
	w, found := c.Roots[rt]
	if !found {
		w = &Word{
			WordID:    c.Max.WordID,
			VByIDX:    map[string]VIDX{},
			Documents: map[DocID]sig{},
		}
		c.Max.WordID++
		c.Roots[rt] = w
		c.IDs[w.WordID] = w
	}
	vidx, found := w.VByIDX[word]
	if !found {
		vidx = VIDX(len(w.Variants))
		w.Variants = append(w.Variants, []byte(word))
		w.VByIDX[word] = vidx
	}
	return w.WordID, vidx
}

// str must start with letterNumber but can have trailing non-letter number
func root(str string) string {
	s := newScanner(str)
	s.matchLetterNumber(false)
	return strings.ToLower(s.str(0, s.idx))
}
