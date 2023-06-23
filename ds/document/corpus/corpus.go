package corpus

import "github.com/adamcolton/luce/ds/document"

const (
	MaxUint32    = ^uint32(0)
	MaxVariantID = VariantID(MaxUint32)
	MaxRootID    = RootID(MaxUint32)
)

// Corpus holds a collection of documents for text indexing. It fulfills
// document.Encoder and document.Decoder.
type Corpus struct {
	Splitter    func(string) (string, []string)
	RootVariant func(str string) (string, document.Variant)

	cur struct {
		RootID
		VariantID
		DocID
	}
	id2root    map[RootID]*root
	rootByStr  map[string]*root
	variant2id map[string]VariantID
	id2variant map[VariantID]document.Variant
	docs       map[DocID]*Document
}

// New creates a Corpus.
func New() *Corpus {
	return &Corpus{
		Splitter:    document.Parse,
		RootVariant: document.RootVariant,
		id2root:     make(map[RootID]*root),
		rootByStr:   make(map[string]*root),
		variant2id:  make(map[string]VariantID),
		id2variant:  make(map[VariantID]document.Variant),
		docs:        make(map[DocID]*Document),
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

// AddDoc to Corpus, returns the string encoded as a Document with a DocID.
func (c *Corpus) AddDoc(str string) *Document {
	enc := &document.DocumentEncoder[RootID, VariantID]{
		Encoder:         c,
		Splitter:        c.Splitter,
		RootVariant:     c.RootVariant,
		WordSingleToken: MaxRootID,
		VarSingleToken:  MaxVariantID,
	}
	d := &Document{
		DocID:    c.cur.DocID,
		Document: enc.Build(str),
		c:        c,
	}
	c.cur.DocID++
	c.docs[d.DocID] = d
	return d
}
