package handler_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	cmds, err := handler.Cmds([]*handler.Command{
		{
			Name: "test",
			Action: func(s string) (int, error) {
				return strconv.Atoi(s)
			},
			Subcmds: []*handler.Command{
				{
					Name: "testSub",
					Action: func() string {
						return "called TestSub"
					},
				},
			},
		}, {
			Name:  "foo",
			Usage: "does foo",
			Action: func(i int) string {
				return strconv.Itoa(i)
			},
			Subcmds: []*handler.Command{
				{
					Name: "SubA",
					Action: func() string {
						return "called SubA"
					},
				}, {
					Name: "SubB",
					Action: func() string {
						return "called SubB"
					},
					Alias: "sb",
				},
			},
		},
	})
	assert.NoError(t, err)

	c, h := cmds.Get([]string{"test"})
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

	c, h = cmds.Get([]string{"foo", "SubA"})
	assert.Equal(t, "SubA", c.Name)
	a, err = h.Handle(nil)
	assert.NoError(t, err)
	assert.Equal(t, "called SubA", a)

	buf := bytes.NewBuffer(nil)
	w := cmds.Writer(nil)
	n, err := w.WriteTo(buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)
	assert.Equal(t, "   foo  does foo\n   test", buf.String())

	buf.Reset()
	w.Recursive = true
	n, err = w.WriteTo(buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)
	assert.Equal(t, "   foo  does foo\n      SubA\n      sb, SubB\n   test\n      testSub", buf.String())

	buf.Reset()
	n, err = cmds.Writer([]string{"foo"}).WriteTo(buf)
	assert.NoError(t, err)
	assert.True(t, n > 0)
	assert.Equal(t, "   SubA\n   sb, SubB", buf.String())

	c, h = cmds.Get([]string{"sb"})
	assert.Equal(t, "SubB", c.Name)
	a, err = h.Handle(nil)
	assert.NoError(t, err)
	assert.Equal(t, "called SubB", a)

}
