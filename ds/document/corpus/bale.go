package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/entity"
)

var (
	doc2keyTransform = morph.OnVal[document.ID]((*docRef).EntKey)
	key2docTransform = morph.OnVal[document.ID](entity.KeyRef[Document])
	variantTransform = morph.NewKVToKV(func(id VariantID, v document.Variant) (string, VariantID) {
		return string(v), id
	})
)

type CorpusBale struct {
	Prefix entity.Key
	Cur    struct {
		RootID
		VariantID
		document.ID
	}
	RootWords []*RootBale
	Variants  map[VariantID]document.Variant
	Docs      map[document.ID]entity.Key
}

func (c *Corpus) TypeID32() uint32 {
	return 2694953506
}

func (c *Corpus) Bale() *CorpusBale {
	return &CorpusBale{
		Prefix:    c.prefix.EntKey(),
		Cur:       c.cur,
		RootWords: c.getRootBales(),
		Variants:  c.id2variant.Copy(0),
		Docs:      lmap.FromEacher(doc2keyTransform.Eacher(c.docs)).Map(),
	}

}

func (bale *CorpusBale) EntRefs() []entity.Key {
	out := make([]entity.Key, 1, len(bale.Docs)+1)
	out[0] = bale.Prefix
	for _, k := range bale.Docs {
		out = append(out, k)
	}
	return out
}

func (bale *CorpusBale) UnbaleTo(c *Corpus) {
	pr := entity.KeyRef[prefix.Prefix](bale.Prefix)
	rLn := len(bale.RootWords)
	c.id2root = lmap.EmptySafe[RootID, *root](rLn)
	c.rootByStr = lmap.EmptySafe[string, *root](rLn)
	for _, rb := range bale.RootWords {
		r := rb.unbale()
		c.id2root.Set(r.RootID, r)
		c.rootByStr.Set(r.str, r)
	}
	variantsBuf := lmap.EmptySafe[string, VariantID](len(bale.Variants))
	c.variant2id = variantsBuf.AppendEacher(variantTransform.Eacher(lmap.New(bale.Variants)))
	docsBuf := lmap.EmptySafe[document.ID, *docRef](len(bale.Docs))
	c.docs = docsBuf.AppendEacher(key2docTransform.Eacher(lmap.New(bale.Docs)))
	c.Splitter = document.Parse
	c.RootVariant = document.RootVariant
	c.Root = document.Root
	c.prefix = pr
	c.cur = bale.Cur
	c.id2variant = lmap.NewSafe(bale.Variants)
}

func (bale *CorpusBale) Unbale() *Corpus {
	out := &Corpus{}
	out.SetDefaults()
	bale.UnbaleTo(out)
	return out
}

func (c *Corpus) getRootBales() []*RootBale {
	roots := make([]*RootBale, 0, c.id2root.Len())
	c.id2root.Each(func(id RootID, r *root, done *bool) {
		roots = append(roots, r.bale())
	})
	return roots
}

type RootBale struct {
	RootID
	Str  string
	Docs []document.ID
}

func (r *root) TypeID32() uint32 {
	return 976961939
}

func (r *root) bale() *RootBale {
	return &RootBale{
		RootID: r.RootID,
		Str:    r.str,
		Docs:   r.docs.Slice(nil),
	}
}

func (bale *RootBale) unbale() *root {
	r := &root{}
	bale.unbaleTo(r)
	return r
}
func (bale *RootBale) unbaleTo(r *root) {
	r.RootID = bale.RootID
	r.str = bale.Str
	r.docs = lset.New(bale.Docs...)
}

type DocumentBale struct {
	DocBale *DocBaleType
	Corpus  entity.Key
}

func (d *Document) TypeID32() uint32 {
	return 1234822691
}

func (d *Document) Bale() *DocumentBale {
	return &DocumentBale{
		DocBale: d.DocType.Bale(),
		Corpus:  d.c.EntKey(),
	}
}

func (bale *DocumentBale) UnbaleTo(d *Document) {
	d.c = entity.KeyRef[Corpus](bale.Corpus)
	d.DocType = &DocType{}
	bale.DocBale.UnbaleTo(d.DocType)
}

func (bale *DocumentBale) EntRefs() []entity.Key {
	return nil
}
