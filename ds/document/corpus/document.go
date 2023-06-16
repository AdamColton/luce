package corpus

import "github.com/adamcolton/luce/ds/document"

type DocID uint32

func (id DocID) DocID() DocID {
	return id
}

type DocIDer interface {
	ID() DocID
}

type Document struct {
	DocID
	*document.Document[RootID, VariantID]
	c *Corpus
}

func (d *Document) String() string {
	dec := &document.DocumentDecoder[RootID, VariantID]{
		Decoder:         d.c,
		WordSingleToken: MaxRootID,
		VarSingleToken:  MaxVariantID,
	}
	return dec.Decode(d.Document)
}

type Documents []*Document

func (ds Documents) Strings() []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		if d != nil {
			out[i] = d.String()
		}
	}
	return out
}
