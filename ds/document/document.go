package document

import (
	"github.com/adamcolton/luce/ds/huffman/huffslice"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
)

// Encoder supplies the necessary encoding information to translate strings
// into IDs.
type Encoder[WordID, VariantID any] interface {
	WordToID(string) WordID
	VariantToID(Variant) VariantID
}

// Encoder supplies the necessary decoding information to translate IDs
// into strings.
type Decoder[WordID, VariantID any] interface {
	IDToWord(WordID) string
	IDToVariant(VariantID) Variant
}

// DocumentEncoder can encode a string into a Document.
type DocumentEncoder[WordID, VariantID comparable] struct {
	Encoder[WordID, VariantID]
	Splitter        func(string) (string, []string)
	RootVariant     func(string) (string, Variant)
	WordSingleToken WordID
	VarSingleToken  VariantID
}

// DocumentDecoder can decode a Document into a string.
type DocumentDecoder[WordID, VariantID comparable] struct {
	Decoder[WordID, VariantID]
	WordSingleToken WordID
	VarSingleToken  VariantID
}

// Document is a string encoded as the root words. This makes identifying
// which words are in a document fast.
type Document[WordID, VariantID comparable] struct {
	Start            string
	ByteLen, WordLen int
	// Words holds the root words present in the document. This slice can
	// be reordered without effecting the encoding.
	Words    []Locations[WordID]
	Variants *huffslice.Slice[VariantID]
}

// Locations hold an ID and the index locations where that ID occures.
type Locations[T comparable] struct {
	ID   T
	Idxs []uint32
}

// WordIDs returns a slice with all the WordIDs in the document.
func (doc *Document[WordID, VariantID]) WordIDs() []WordID {
	ln := len(doc.Words)

	out := make([]WordID, 0, ln)
	for _, w := range doc.Words {
		out = append(out, w.ID)
	}
	return out
}

// Build takes a stirng and encodes it to a Document.
func (enc DocumentEncoder[WordID, VariantID]) Build(str string) *Document[WordID, VariantID] {
	doc := &Document[WordID, VariantID]{}
	enc.build(doc, str)
	return doc
}

func (enc DocumentEncoder[WordID, VariantID]) build(doc *Document[WordID, VariantID], str string) {
	start, words := enc.Splitter(str)
	doc.Start = start
	doc.ByteLen = len(str)
	doc.WordLen = len(words)

	wls := make(map[WordID]Locations[WordID])
	vEnc := huffslice.NewEncoder(doc.WordLen, enc.VarSingleToken)

	for idx, w := range words {
		root, variant := enc.RootVariant(w)
		wID := enc.WordToID(root)
		vID := enc.VariantToID(variant)

		wl := wls[wID]
		wl.ID = wID
		wl.Idxs = append(wl.Idxs, uint32(idx))
		wls[wID] = wl

		vEnc.Slice = append(vEnc.Slice, vID)
	}

	doc.Words = make([]Locations[WordID], 0, len(wls))
	doc.Variants = vEnc.Encode()

	for _, wl := range wls {
		doc.Words = append(doc.Words, wl)
	}
}

// Decode a Document to a string
func (dec DocumentDecoder[WordID, VariantID]) Decode(doc *Document[WordID, VariantID]) string {
	words := slice.Make[string](doc.WordLen, 0)
	out := make([]byte, len(doc.Start), doc.ByteLen)
	copy(out, []byte(doc.Start))

	for _, wl := range doc.Words {
		w := dec.IDToWord(wl.ID)
		for _, idx := range wl.Idxs {
			words[idx] = w
		}

	}

	vi := doc.Variants.Iter()
	for vID, done := vi.Cur(); !done; vID, done = vi.Next() {
		out = dec.IDToVariant(vID).Apply(words[vi.Idx()], out)
	}
	return string(out)
}

// ChangeSet shows the IDs that were added and removed from a Document.
type ChangeSet[T any] struct {
	Added, Removed slice.Slice[T]
}

// Update a document updates the encoding and returns a ChangeSet.
func (enc DocumentEncoder[WordID, VariantID]) Update(doc *Document[WordID, VariantID], str string) *ChangeSet[WordID] {
	wordsBefore := lset.New(doc.WordIDs()...)
	enc.build(doc, str)
	cs := &ChangeSet[WordID]{}
	for _, wID := range doc.WordIDs() {
		if wordsBefore.Contains(wID) {
			wordsBefore.Remove(wID)
		} else {
			cs.Added = append(cs.Added, wID)
		}
	}
	cs.Removed = wordsBefore.Slice(nil)
	return cs
}
