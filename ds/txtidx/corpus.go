package txtidx

import (
	"strings"

	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/util/lstr"
)

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	roots         *prefix.Prefix
	words         []*word
	wordsByStr    map[string]*word
	variantsByStr map[string]varIDX
	variants      []variant
	docs          []*document
	unused        struct {
		docs  []DocID
		words []wordIDX
	}
}

func NewCorpus() *Corpus {
	return &Corpus{
		roots:         prefix.New(),
		wordsByStr:    make(map[string]*word),
		variantsByStr: make(map[string]varIDX),
	}
}

type sig struct{}

func (c *Corpus) upsert(wrd string) (wordIDX, varIDX) {
	rt := root(wrd)
	w := c.wordsByStr[rt]
	if w == nil {
		w = &word{
			str:       rt,
			Documents: newDocSet(),
		}
		c.roots.Upsert(rt)
		c.wordsByStr[rt] = w
		ln := len(c.unused.words)
		if ln > 0 {
			ln--
			w.wordIDX = c.unused.words[ln]
			c.unused.words = c.unused.words[:ln]
			c.words[w.wordIDX] = w
		} else {
			w.wordIDX = wordIDX(len(c.words))
			c.words = append(c.words, w)
		}
	}
	v := findVariant(rt, wrd)
	vid, found := c.variantsByStr[string(v)]
	if !found {
		vid = varIDX(len(c.variants))
		c.variantsByStr[string(v)] = vid
		c.variants = append(c.variants, v)
	}
	return w.wordIDX, vid
}

func (c *Corpus) Find(words ...string) DocSet {
	return c.find(words...)
}

func (c *Corpus) find(terms ...string) *docSet {
	if len(terms) == 0 {
		return newDocSet()
	}

	out := c.findSingle(terms[0])
	for _, w := range terms[1:] {
		if out == nil {
			break
		}
		out.intersectMerge(c.findSingle(w))
	}
	return out
}

func (c *Corpus) findSingle(word string) *docSet {
	var ws words
	for _, n := range c.roots.Containing(root(word)) {
		for _, w := range n.AllWords() {
			g := w.Gram()
			w := c.wordsByStr[g]
			if w == nil {
				panic("expected word")
			}
			ws = append(ws, w)
		}
	}
	return ws.docSetUnion()
}

func (c *Corpus) AddDoc(doc string) Document {
	return newDoc(doc, c)
}

func (c *Corpus) allocDocIDX(d *document) {
	ln := len(c.unused.docs)
	if ln > 0 {
		ln--
		d.id = c.unused.docs[ln]
		c.unused.docs = c.unused.docs[:ln]
		c.docs[d.id] = d
	} else {
		d.id = DocID(len(c.docs))
		c.docs = append(c.docs, d)
	}
}

func (c *Corpus) Delete(di DocIDer) {
	d := c.getDoc(di)
	id := di.ID()
	c.docs[id] = nil
	c.unused.docs = append(c.unused.docs, id)

	for _, wIdx := range d.words() {
		c.deleteDocWord(id, c.words[wIdx])
	}
	c.roots.Purge()
}

func (c *Corpus) deleteDocWord(di DocIDer, w *word) {
	w.Documents.Delete(di)
	if w.Documents.Len() == 0 {
		c.words[w.wordIDX] = nil
		delete(c.wordsByStr, w.str)
		c.roots.Remove(w.str)
		c.unused.words = append(c.unused.words, w.wordIDX)
	}
}

func (c *Corpus) getDoc(id DocIDer) *document {
	d, ok := id.(*document)
	if !ok {
		d = c.docs[id.ID()]
	}
	return d
}

func (c *Corpus) Update(id DocIDer, txt string) {
	c.getDoc(id).update(c, txt)
}

func (c *Corpus) Search(search string) (DocSet, []string) {

	s := buildSearch(lstr.NewScanner(search))
	ds := c.find(s.words...).copy()
	var strs []string
	if len(s.exact) > 0 {
		for _, di := range ds.t.Slice() {
			str := c.docs[di].toString(c)
			for _, e := range s.exact {
				if !strings.Contains(str, e) {
					ds.Delete(di)
					break
				}
			}
			if ds.Has(di) {
				strs = append(strs, str)
			}
		}
	} else {
		strs = c.GetDocs(ds)
	}
	return ds, strs
}

func (c *Corpus) GetDocs(docs DocSet) []string {
	ds := docs.(*docSet)
	out := make([]string, 0, ds.Len())
	for _, di := range ds.t.Slice() {
		out = append(out, c.docs[di].toString(c))
	}
	return out
}

func (c *Corpus) DocString(id DocIDer) string {
	return c.getDoc(id).toString(c)
}

func (c *Corpus) Suggest(word string, max int) []prefix.Suggestion {
	return c.roots.Find(word).Suggest(max)
}
