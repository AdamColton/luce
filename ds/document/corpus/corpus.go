package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/prefix"
)

type VariantID uint32

const (
	MaxUint32    = ^uint32(0)
	MaxVariantID = VariantID(MaxUint32)
	MaxRootID    = RootID(MaxUint32)
)

type Corpus struct {
	Splitter    func(string) (string, []string)
	RootVariant func(str string) (string, document.Variant)
	Root        func(str string) string

	prefix *prefix.Prefix
	cur    struct {
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

func New() *Corpus {
	return &Corpus{
		Splitter:    document.Parse,
		RootVariant: document.RootVariant,
		Root:        document.Root,

		prefix:     prefix.New(),
		id2root:    make(map[RootID]*root),
		rootByStr:  make(map[string]*root),
		variant2id: make(map[string]VariantID),
		id2variant: make(map[VariantID]document.Variant),
		docs:       make(map[DocID]*Document),
	}
}

func (c *Corpus) WordToID(rStr string) RootID {
	r := c.rootByStr[rStr]
	if r == nil {
		r = &root{
			RootID: c.cur.RootID,
			str:    rStr,
			docs:   lset.New[DocID](),
		}
		c.rootByStr[rStr] = r
		c.id2root[c.cur.RootID] = r
		c.cur.RootID++
		c.prefix.Upsert(rStr)
	}
	return r.RootID
}

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

func (c *Corpus) IDToWord(rID RootID) string {
	r := c.id2root[rID]
	if r == nil {
		return ""
	}
	return r.str
}

func (c *Corpus) IDToVariant(vID VariantID) document.Variant {
	return c.id2variant[vID]
}

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

	for _, rID := range d.WordIDs() {
		c.id2root[rID].docs.Add(d.DocID)
	}
	return d
}

func (c *Corpus) Find(word string) *lset.Set[DocID] {
	r := c.rootByStr[c.Root(word)]
	if r == nil {
		return nil
	}
	return r.docs
}

func (c *Corpus) Prefix(gram string) prefix.Node {
	return c.prefix.Find(gram)
}

func (c *Corpus) Containing(gram string) prefix.Nodes {
	return c.prefix.Containing(gram)
}

func (c *Corpus) GetDoc(id DocID) *Document {
	return c.docs[id]
}

func (c *Corpus) GetDocs(ids []DocID) Documents {
	out := make(Documents, len(ids))
	for i, id := range ids {
		out[i] = c.docs[id]
	}
	return out
}
