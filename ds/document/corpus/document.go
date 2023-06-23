package corpus

import "github.com/adamcolton/luce/ds/document"

// DocID allows references to a document to be passed around
type DocID uint32

// ID fullfils DocIDer
func (id DocID) ID() DocID {
	return id
}

// DocIDer allows anything that can reference a DocID to be used to retreive
// a document.
type DocIDer interface {
	ID() DocID
}

// Document uses a Corpus to fulfill Encoder and Decoder for document.Document.
type Document struct {
	DocID
	*document.Document[RootID, VariantID]
	c *Corpus
}

// String decodes a document
func (d *Document) String() string {
	dec := &document.DocumentDecoder[RootID, VariantID]{
		Decoder:         d.c,
		WordSingleToken: MaxRootID,
		VarSingleToken:  MaxVariantID,
	}
	return dec.Decode(d.Document)
}
