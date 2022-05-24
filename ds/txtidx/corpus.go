package txtidx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	Roots         *Markov
	Words         []*Word
	Variants      map[string]VarIDX
	VariantLookup []variant
	Docs          []*Document
	unused        struct {
		docs  []DocIDX
		words []WordIDX
	}
}

func NewCorpus() *Corpus {
	return &Corpus{
		Roots:    NewMarkov(),
		Variants: map[string]VarIDX{},
	}
}

type sig struct{}

type Word struct {
	Word string
	WordIDX
	Documents *DocSet
}

type WordIDX uint32
type VarIDX uint32

type Suffix []byte

func (c *Corpus) Get(word string) (WordIDX, VarIDX) {
	wid, vid := WordIDX(MaxUint32), VarIDX(MaxUint32)
	rt := root(word)
	v := findVariant(rt, word)
	tmpVid, found := c.Variants[string(v)]
	if found {
		vid = tmpVid
	}

	w := c.Roots.Find(rt)
	if w != nil {
		wid = w.WordIDX
	}

	return wid, vid
}

func (c *Corpus) Upsert(word string) (WordIDX, VarIDX) {
	rt := root(word)
	w := c.Roots.Upsert(rt)
	if w.WordIDX == WordIDX(MaxUint32) {
		ln := len(c.unused.words)
		if ln > 0 {
			ln--
			w.WordIDX = c.unused.words[ln]
			c.unused.words = c.unused.words[:ln]
			c.Words[w.WordIDX] = w
		} else {
			w.WordIDX = WordIDX(len(c.Words))
			c.Words = append(c.Words, w)
		}
		w.Word = rt
	}
	v := findVariant(rt, word)
	vid, found := c.Variants[string(v)]
	if !found {
		vid = VarIDX(len(c.VariantLookup))
		c.Variants[string(v)] = vid
		c.VariantLookup = append(c.VariantLookup, v)
	}
	return w.WordIDX, vid
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

func (c *Corpus) Find(words ...string) *DocSet {
	if len(words) == 0 {
		return newDocSet()
	}
	out := c.find(words[0])

	for _, w := range words[1:] {
		if out == nil {
			break
		}
		out = out.Intersect(c.find(w))
	}
	return out
}

func (c *Corpus) find(word string) *DocSet {
	ws := c.Roots.FindAll(root(word))
	if len(ws) == 0 {
		return newDocSet()
	}
	out := ws[0].Documents.Copy()
	for _, w := range ws[1:] {
		out.Merge(w.Documents)
	}
	return out
}

func (c *Corpus) AddDoc(doc string) *Document {
	pp := c.newPP()
	pp.set(doc)
	return pp.build()
}

func (c *Corpus) allocDocIDX(d *Document) {
	ln := len(c.unused.docs)
	if ln > 0 {
		ln--
		d.DocIDX = c.unused.docs[ln]
		c.unused.docs = c.unused.docs[:ln]
		c.Docs[d.DocIDX] = d
	} else {
		d.DocIDX = DocIDX(len(c.Docs))
		c.Docs = append(c.Docs, d)
	}
}

func (c *Corpus) Delete(di DocIDX) {
	d := c.Docs[di]
	c.Docs[di] = nil
	c.unused.docs = append(c.unused.docs, di)

	for _, dw := range d.Words {
		w := c.Words[dw.WordIDX]
		w.Documents.Delete(di)
		if w.Documents.Len() == 0 {
			c.deleteWord(w)
			c.unused.words = append(c.unused.words, w.WordIDX)
		}
	}
}

func (c *Corpus) deleteWord(w *Word) {
	c.Words[w.WordIDX] = nil
	c.Roots.deleteWord(w.Word)
}
