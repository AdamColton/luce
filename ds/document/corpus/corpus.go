package corpus

import "github.com/adamcolton/luce/ds/document"

// Corpus holds a collection of documents for text indexing. It fulfills
// document.Encoder and document.Decoder.
type Corpus struct {
	cur struct {
		RootID
		VariantID
	}
	id2root    map[RootID]*root
	rootByStr  map[string]*root
	variant2id map[string]VariantID
	id2variant map[VariantID]document.Variant
}

// New creates a Corpus.
func New() *Corpus {
	return &Corpus{
		id2root:    make(map[RootID]*root),
		rootByStr:  make(map[string]*root),
		variant2id: make(map[string]VariantID),
		id2variant: make(map[VariantID]document.Variant),
	}
}

// WordToID converts a root word to a RootID, fulfilling document.Encoder.
func (c *Corpus) WordToID(rStr string) RootID {
	r := c.rootByStr[rStr]
	if r == nil {
		r = &root{
			RootID: c.cur.RootID,
			str:    rStr,
		}
		c.rootByStr[rStr] = r
		c.id2root[c.cur.RootID] = r
		c.cur.RootID++
	}
	return r.RootID
}

// VariantToID converts a root word to a VariantID, fulfilling document.Encoder.
func (c *Corpus) VariantToID(v document.Variant) VariantID {
	vID, found := c.variant2id[string(v)]
	if !found {
		vID = c.cur.VariantID
		c.cur.VariantID++
		c.variant2id[string(v)] = vID
		c.id2variant[vID] = v
	}

	return vID
}

// IDToWord converts a RootID to a root word, fulfilling document.Decoder.
func (c *Corpus) IDToWord(rID RootID) string {
	r := c.id2root[rID]
	if r == nil {
		return ""
	}
	return r.str
}

// IDToVariant converts a VariantID to a document.Variant, fulfilling document.Decoder.
func (c *Corpus) IDToVariant(vID VariantID) document.Variant {
	return c.id2variant[vID]
}
