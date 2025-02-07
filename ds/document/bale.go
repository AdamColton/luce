package document

import (
	"github.com/adamcolton/luce/ds/bimap"
	"github.com/adamcolton/luce/ds/huffman/huffslice"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/util/reflector"
)

type DocumentBale[WordID, VariantID comparable] struct {
	ID
	Start            string
	ByteLen, WordLen int
	// Words holds the root words present in the document. This slice can
	// be reordered without effecting the encoding.
	Words    []Locations[WordID]
	Variants *huffslice.SliceBale[VariantID]
}

var typeIDs = bimap.New[uint32, docTypeKey](0)

func AddTypeID[WordID, VariantID comparable](id uint32) {
	k := docTypeKey{
		WordID:    reflector.Type[WordID](),
		VariantID: reflector.Type[VariantID](),
	}
	typeIDs.Add(id, k)
}

func (enc DocumentEncoder[WordID, VariantID]) AddTypeID(id uint32) {
	AddTypeID[WordID, VariantID](id)
}

func (bale *DocumentBale[WordID, VariantID]) TypeID32() uint32 {
	k := docTypeKey{
		WordID:    reflector.Type[WordID](),
		VariantID: reflector.Type[VariantID](),
	}
	id, _ := typeIDs.B(k)
	return id
}

func (bale *DocumentBale[WordID, VariantID]) EntRefs() []entity.Key {
	return nil
}

func (doc *Document[WordID, VariantID]) Bale() *DocumentBale[WordID, VariantID] {
	return &DocumentBale[WordID, VariantID]{
		ID:       doc.ID,
		Start:    doc.Start,
		ByteLen:  doc.ByteLen,
		WordLen:  doc.WordLen,
		Words:    doc.Words,
		Variants: doc.Variants.Bale(),
	}
}

func (bale *DocumentBale[WordID, VariantID]) UnbaleTo(doc *Document[WordID, VariantID]) {
	doc.ID = bale.ID
	doc.Start = bale.Start
	doc.ByteLen = bale.ByteLen
	doc.WordLen = bale.WordLen
	doc.Words = bale.Words
	doc.Variants = bale.Variants.Unbale()
}
