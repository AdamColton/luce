package handler_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	cmds, err := handler.Cmds([]handler.Command{
		{
			Name: "test",
			Action: func(s string) (int, error) {
				return strconv.Atoi(s)
			},
		}, {
			Name: "foo",
			Action: func(i int) string {
				return strconv.Itoa(i)
			},
		},
	})
	assert.NoError(t, err)

	c, h := cmds.Get("test")
	assert.Equal(t, "test", c.Name)
	assert.Equal(t, ltype.String, h.Type())

	a, err := h.Handle("5")
	assert.NoError(t, err)
	assert.Equal(t, 5, a)

	s, err := cmds.Switch()
	assert.NoError(t, err)

	a, err = s.Handle("6")
	assert.NoError(t, err)
	assert.Equal(t, 6, a)

	a, err = s.Handle(52)
	assert.NoError(t, err)
	assert.Equal(t, "52", a)
}
