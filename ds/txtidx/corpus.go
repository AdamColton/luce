package txtidx

const MaxUint32 uint32 = ^uint32(0)

type Corpus struct {
	roots         *markov
	words         []*word
	variantsByStr map[string]varIDX
	variants      []variant
	unused        struct {
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

func (c *Corpus) Suggest(word string, max int) []string {
	return c.roots.suggest(word, max)
}
