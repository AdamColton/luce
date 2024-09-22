package handler_test

import (
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

func TestAddAlias(t *testing.T) {
	// The Commands from a registrar are an lmap.Wrapper, which AddAlias takes.
	cmds := handler.DefaultRegistrar.Commands(handlerObj{name: "test"})
	assert.Equal(t, 2, cmds.Len())

	handler.AddAlias(cmds,
		"string", "s",
		"missing", "m",
		"int", "i",
		"dangling",
	)
	assert.Equal(t, "s", cmds.GetVal("string").Alias)
	assert.Equal(t, "i", cmds.GetVal("int").Alias)
	assert.Equal(t, 2, cmds.Len())
}
