package txtidx

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	roots         *markov
	words         []*word
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
		roots:         newMarkov(),
		variantsByStr: map[string]varIDX{},
	}
}

type sig struct{}

func (c *Corpus) upsert(word string) (wordIDX, varIDX) {
	rt := root(word)
	w := c.roots.upsert(rt)
	if w.wordIDX == wordIDX(MaxUint32) {
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
		w.str = rt
	}
	v := findVariant(rt, word)
	vid, found := c.variantsByStr[string(v)]
	if !found {
		vid = varIDX(len(c.variants))
		c.variantsByStr[string(v)] = vid
		c.variants = append(c.variants, v)
	}
	return w.wordIDX, vid
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
}

func (c *Corpus) deleteDocWord(di DocIDer, w *word) {
	w.Documents.Delete(di)
	if w.Documents.Len() == 0 {
		c.words[w.wordIDX] = nil
		c.roots.deleteWord(w.str)
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

func (c *Corpus) GetDocs(docs DocSet) []string {
	ds := docs.(*docSet)
	out := make([]string, 0, ds.Len())
	for di := range ds.docs {
		out = append(out, c.docs[di].toString(c))
	}
	return out
}

func (c *Corpus) DocString(id DocIDer) string {
	return c.getDoc(id).toString(c)
}
