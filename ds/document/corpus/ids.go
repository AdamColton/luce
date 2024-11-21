package corpus

import "github.com/adamcolton/luce/ds/lset"

// RootID represents a root word - a lower case alphanumeric string
type RootID uint32

type root struct {
	RootID
	str  string
	docs *lset.Set[DocID]
}

// VariantID represents a Variant by ID.
type VariantID uint32
