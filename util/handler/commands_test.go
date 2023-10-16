package handler_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	cmds, err := handler.Cmds([]handler.Command{
		{
			Name: "test",
			Action: func() (int, error) {
				return 5, nil
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
	assert.Equal(t, "test", h.Name())
	assert.Equal(t, "test", c.Name)
	assert.Equal(t, nil, h.Type())

	a, err := h.Handle(nil)
	assert.NoError(t, err)
	assert.Equal(t, 5, a)

	s, err := cmds.Switch()
	assert.NoError(t, err)

	a, err = s.Handle("test")
	assert.NoError(t, err)
	assert.Equal(t, 5, a)

	a, err = s.Handle(52)
	assert.NoError(t, err)
	assert.Equal(t, "52", a)
}
