package txtidx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	Roots         *Markov
	IDs           map[WordID]*Word
	Variants      map[string]VarID
	VariantLookup map[VarID]variant
	Max           struct {
		WordID
		DocID
		VarID
	}
}

func NewCorpus() *Corpus {
	return &Corpus{
		Roots:         NewMarkov(),
		IDs:           map[WordID]*Word{},
		Variants:      map[string]VarID{},
		VariantLookup: map[VarID]variant{},
	}
}

type sig struct{}

type Word struct {
	Word string
	WordID
	Documents *DocSet
}

type WordID uint32
type VarID uint32
type VIDX uint32

type DocRef struct {
	VarRefs []VarRef
}

type VarRef struct {
	VIDX
	Locs []uint32
}

type Suffix []byte

func (c *Corpus) Get(word string) (WordID, VarID) {
	wid, vid := WordID(MaxUint32), VarID(MaxUint32)
	rt := root(word)
	v := findVariant(rt, word)
	tmpVid, found := c.Variants[string(v)]
	if found {
		vid = tmpVid
	}

	w := c.Roots.Find(rt)
	if w != nil {
		wid = w.WordID
	}

	return wid, vid
}

func (c *Corpus) Upsert(word string) (WordID, VarID) {
	rt := root(word)
	w := c.Roots.Upsert(rt)
	if w.WordID == WordID(MaxUint32) {
		w.WordID = c.Max.WordID
		w.Word = rt
		c.Max.WordID++
		c.IDs[w.WordID] = w
	}
	v := findVariant(rt, word)
	vid, found := c.Variants[string(v)]
	if !found {
		vid = c.Max.VarID
		c.Max.VarID++
		c.Variants[string(v)] = vid
		c.VariantLookup[vid] = v
	}
	return w.WordID, vid
}

// str must start with letterNumber but can have trailing non-letter number
func root(str string) string {
	s := newScanner(str)
	s.matchLetterNumber(false)
	return strings.ToLower(s.str(0, s.idx))
}

// divUp division `round up
func divUp(a, b int) int {
	out := a / b
	if out*b != a {
		out++
	}
	return out
}

type variant []byte

func findVariant(root, str string) variant {
	rs := []rune(root)
	b := []byte(str)
	suffix := str[len(root):]
	out := make([]byte, 0, divUp(len(rs), 8)+len(suffix))

	bIdx := 0
	caseByte := byte(0)
	for _, rr := range rs {
		r, ln := utf8.DecodeRune(b)
		b = b[ln:]
		caseByte <<= 1
		if r != rr {
			caseByte |= 1
		}
		bIdx++
		if bIdx == 8 {
			out = append(out, caseByte)
			bIdx = 0
			caseByte = 0
		}
	}
	if bIdx != 0 {
		caseByte <<= 8 - byte(bIdx)
		out = append(out, caseByte)
	}
	out = append(out, suffix...)
	return out
}

const startMask byte = 128

func (v variant) apply(rt string) string {
	ln := len(rt)
	out := make([]byte, 0, len(v)-divUp(ln, 8)+len(rt))

	var mask, caseByte byte
	bIdx := 0
	in := []byte(rt)
	for len(in) > 0 {
		if mask == 0 {
			mask = startMask
			caseByte = v[bIdx]
			bIdx++
		}
		r, size := utf8.DecodeRune(in)
		in = in[size:]
		if caseByte&mask != 0 {
			r = unicode.ToUpper(r)
		}
		mask >>= 1

		out = append(out, string(r)...)

	}

	out = append(out, v[bIdx:]...)
	return string(out)
}

func (c *Corpus) Find(word string) *DocSet {
	w := c.Roots.Find(word)
	if w == nil {
		return nil
	}

	return w.Documents
}
