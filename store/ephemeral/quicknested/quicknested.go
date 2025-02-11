package quicknested

import (
	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/store/flatwrap"
)

func New(bufferSize int) store.NestedFactory {
	root := ephemeral.Factory(bytebtree.New, bufferSize)
	return flatwrap.New(root)
}
