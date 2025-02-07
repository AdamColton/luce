package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/entity"
)

// RootID represents a root word - a lower case alphanumeric string
type RootID uint32

type root struct {
	entity.Key
	RootID
	str  string
	docs *lset.Set[document.ID]
}

// VariantID represents a Variant by ID.
type VariantID uint32
