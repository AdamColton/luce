package byteid

import "github.com/adamcolton/luce/ds/idx"

// Index allows the equivalent of map[[]byte]<Type>.
type Index = idx.Index[[]byte]

type IndexFactory = idx.IndexFactory[[]byte]
