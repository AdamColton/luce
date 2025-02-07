package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/math/ints"
)

const (
	MaxVariantID = VariantID(ints.MaxU32)
	MaxRootID    = RootID(ints.MaxU32)
)

type prefixRef = entity.Ref[prefix.Prefix, *prefix.Prefix]
type docRef = entity.Ref[Document, *Document]

// Corpus holds a collection of documents for text indexing. It fulfills
// document.Encoder and document.Decoder.
type Corpus struct {
	Splitter    func(string) (string, []string)
	RootVariant func(str string) (string, document.Variant)
	Root        func(str string) string

	key    entity.Key
	prefix *prefixRef
	cur    struct {
		RootID
		VariantID
		document.ID
	}
	id2root    lmap.Wrapper[RootID, *root]
	rootByStr  lmap.Wrapper[string, *root]
	variant2id lmap.Wrapper[string, VariantID]
	id2variant lmap.Wrapper[VariantID, document.Variant]
	docs       lmap.Wrapper[document.ID, *docRef]
	save       bool
	ref        *entity.Ref[Corpus, *Corpus]
}

func NewKey(key entity.Key) *Corpus {
	p := prefix.New()
	pr := entity.Put(p)
	c := &Corpus{
		key:        key,
		prefix:     pr,
		id2root:    lmap.EmptySafe[RootID, *root](0),
		rootByStr:  lmap.EmptySafe[string, *root](0),
		variant2id: lmap.EmptySafe[string, VariantID](0),
		id2variant: lmap.EmptySafe[VariantID, document.Variant](0),
		docs:       lmap.EmptySafe[document.ID, *docRef](0),
	}
	c.ref = entity.Put(c)
	c.SetDefaults()
	return c
}

// New creates a Corpus.
func New() *Corpus {
	return NewKey(entity.Rand())
}

var (
	DefaultSpitter     = document.Parse
	DefaultRootVariant = document.RootVariant
	DefaultRoot        = document.Root
)

func (c *Corpus) SetDefaults() {
	c.Splitter = DefaultSpitter
	c.RootVariant = DefaultRootVariant
	c.Root = DefaultRoot
}

func (c *Corpus) getRootByStr(rStr string) *root {
	r, ok := c.rootByStr.Get(rStr)
	if !ok {
		return nil
	}
	return r
}

// WordToID converts a root word to a RootID, fulfilling document.Encoder.
func (c *Corpus) WordToID(rStr string) RootID {
	r := c.getRootByStr(rStr)
	if r == nil {
		r = &root{
			Key:    entity.Rand(),
			RootID: c.cur.RootID,
			str:    rStr,
			docs:   lset.New[document.ID](),
		}
		c.rootByStr.Set(rStr, r)
		c.id2root.Set(c.cur.RootID, r)
		c.cur.RootID++
		c.prefix.GetPtr().Upsert(rStr)
		c.saveIf(c)
	}
	return r.RootID
}

// VariantToID converts a root word to a VariantID, fulfilling document.Encoder.
func (c *Corpus) VariantToID(v document.Variant) VariantID {
	vID, found := c.variant2id.Get(string(v))
	if !found {
		vID = c.cur.VariantID
		c.cur.VariantID++
		c.variant2id.Set(string(v), vID)
		c.id2variant.Set(vID, v)
		c.saveIf(c)
	}

	return vID
}

func (c *Corpus) getRoot(rID RootID) *root {
	r, ok := c.id2root.Get(rID)
	if !ok {
		return nil
	}
	return r
}

// IDToWord converts a RootID to a root word, fulfilling document.Decoder.
func (c *Corpus) IDToWord(rID RootID) string {
	r := c.getRoot(rID)
	if r == nil {
		return ""
	}
	return r.str
}

// IDToVariant converts a VariantID to a document.Variant, fulfilling document.Decoder.
func (c *Corpus) IDToVariant(vID VariantID) document.Variant {
	return c.id2variant.GetVal(vID)
}

func (c *Corpus) Encoder() *document.DocumentEncoder[RootID, VariantID] {
	return &document.DocumentEncoder[RootID, VariantID]{
		Encoder:         c,
		Splitter:        c.Splitter,
		RootVariant:     c.RootVariant,
		WordSingleToken: MaxRootID,
		VarSingleToken:  MaxVariantID,
	}
}

// AddDoc to Corpus, returns the string encoded as a Document with a DocID.
func (c *Corpus) AddDoc(str string) *Document {
	enc := c.Encoder()
	id := c.cur.ID
	d := &Document{
		DocType: enc.Build(str),
		c:       c.ref,
	}
	d.DocType.ID = id
	c.cur.ID++
	c.docs.Set(id, entity.Put(d))

	for _, rID := range d.WordIDs() {
		r := c.getRoot(rID)
		r.docs.Add(id)
	}
	c.saveIf(c, d)

	return d
}

// Find all documents containing a word
func (c *Corpus) Find(word string) *lset.Set[document.ID] {
	r := c.getRootByStr(c.Root(word))
	if r == nil {
		return nil
	}
	return r.docs
}

// Prefix returns the prefix.Node for all words in the corpus.
func (c *Corpus) Prefix(gram string) prefix.Node {
	return c.prefix.GetPtr().Find(gram)
}

// Containing returns a prefix.Nodes for all nodes containing the given gram.
func (c *Corpus) Containing(gram string) prefix.Nodes {
	return c.prefix.GetPtr().Containing(gram)
}

// GetDoc returns a Document by DocID.
func (c *Corpus) GetDoc(id document.ID) *Document {
	er, ok := c.docs.Get(id)
	if !ok {
		return nil
	}
	d, ok := er.Get()
	if !ok {
		return nil
	}
	return d
}

func (c *Corpus) AllIDs() []document.ID {
	return c.docs.Keys(nil)
}

// GetDocs returns a set of documents by DocID
func (c *Corpus) GetDocs(ids []document.ID) Documents {
	out := make(Documents, len(ids))
	for i, id := range ids {
		out[i] = c.GetDoc(id)
	}
	return out
}

func (c *Corpus) Remove(id document.ID) {
	doc := c.GetDoc(id)
	if doc == nil {
		return
	}
	for _, l := range doc.Words {
		r := c.getRoot(l.ID)
		r.docs.Remove(id)
	}
	if c.save {
		ref := c.docs.GetVal(id)
		ref.Delete()
	}
	c.docs.Delete(id)
	c.saveIf(c)
}
