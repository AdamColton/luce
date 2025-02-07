package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/entity"
)

// DocIDer allows anything that can reference a DocID to be used to retreive
// a document.
type DocIDer interface {
	DocID() document.ID
}

type DocType = document.Document[RootID, VariantID]
type DocBaleType = document.DocumentBale[RootID, VariantID]

// Document uses a Corpus to fulfill Encoder and Decoder for document.Document.
type Document struct {
	*DocType
	c *entity.Ref[Corpus, *Corpus]
}

// String decodes a document
func (d *Document) String() string {
	c := d.c.GetPtr()
	dec := &document.DocumentDecoder[RootID, VariantID]{
		Decoder:         c,
		WordSingleToken: MaxRootID,
		VarSingleToken:  MaxVariantID,
	}
	return dec.Decode(d.DocType)
}

func (d *Document) entSave() {
	entity.Save(d)
}

func (d *Document) Update(str string) {
	c := d.c.GetPtr()
	cs := c.Encoder().Update(d.DocType, str)
	for _, rm := range cs.Removed {
		r := c.getRoot(rm)
		r.docs.Remove(d.ID)
	}
	for _, add := range cs.Added {
		r := c.getRoot(add)
		r.docs.Add(d.ID)
	}
	c.saveIf(c)
}

// Documents is a collection of documents
type Documents []*Document

// Strings converts all the documents in the collection to stirngs.
func (ds Documents) Strings() []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		if d != nil {
			out[i] = d.String()
		}
	}
	return out
}
