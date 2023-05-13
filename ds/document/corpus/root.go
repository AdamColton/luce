package corpus

import "github.com/adamcolton/luce/ds/lset"

type RootID uint32

type root struct {
	RootID
	str  string
	docs *lset.Set[DocID]
}
