// Package quicknested provides an in-memory nested store that is ready to use.
package quicknested

import (
	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/store/flatwrap"
)

// New creates a store.NestedFactory that keeps everything in memory. Nothing is
// persisted. It is a flatwrap over an ephemeral store, so nested stores are
// buckets in a flat store. bufferSize is the number of records to make room for
// in each store, as it is for ephemeral.Factory. It is not safe for concurrent
// use.
func New(bufferSize int) store.NestedFactory {
	root := ephemeral.Factory(bytebtree.New, bufferSize)
	return flatwrap.New(root)
}
