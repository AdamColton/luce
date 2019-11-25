package byteid

import "github.com/adamcolton/luce/ds/idx"

// Index allows the equivalent of map[[]byte]<Type>. It is idx.Index with
// []byte keys.
type Index = idx.Index[[]byte]

// IndexFactory creates an empty Index for a slice that has a length of
// slicelen. It is idx.IndexFactory with []byte keys.
type IndexFactory = idx.IndexFactory[[]byte]
