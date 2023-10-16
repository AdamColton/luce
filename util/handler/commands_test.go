package handler_test

import (
	"bytes"
	"sort"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

// newCmds is a tree with an alias, a usage, commands with and without an
// argument and the empty-named command that is used when nothing matches.
func newCmds(t *testing.T) *handler.Commands {
	cmds, err := handler.Cmds([]*handler.Command{
		{
			Name:  "foo",
			Usage: "does foo",
			Alias: "f",
			Action: func(i int) string {
				return strconv.Itoa(i)
			},
			Subcmds: []*handler.Command{
				{Name: "SubA", Action: func() string { return "called SubA" }},
				{Name: "SubB", Alias: "sb", Action: func() string { return "called SubB" }},
			},
		}, {
			Name:   "test",
			Action: func(s string) (int, error) { return strconv.Atoi(s) },
		}, {
			Name:   "",
			Action: func() string { return "nothing matched" },
		},
	})
	assert.NoError(t, err)
	return cmds
}

func TestCmdsErrors(t *testing.T) {
	ok := func() string { return "ok" }
	tt := map[string]struct {
		cmds []*handler.Command
		err  string
	}{
		"bad action": {
			cmds: []*handler.Command{{Name: "bad", Action: "not a func"}},
			err:  "expected func, got: string",
		},
		"bad sub action": {
			cmds: []*handler.Command{{Name: "a", Action: ok, Subcmds: []*handler.Command{{Name: "bad", Action: 5}}}},
			err:  "expected func, got: int",
		},
		"duplicate name": {
			cmds: []*handler.Command{{Name: "a", Action: ok}, {Name: "a", Action: ok}},
			err:  "duplicate command name: a",
		},
		"duplicate sub name": {
			cmds: []*handler.Command{{Name: "a", Action: ok, Subcmds: []*handler.Command{
				{Name: "b", Action: ok}, {Name: "b", Action: ok},
			}}},
			err: "duplicate command name: b",
		},
		"duplicate alias": {
			cmds: []*handler.Command{{Name: "a", Alias: "x", Action: ok}, {Name: "b", Alias: "x", Action: ok}},
			err:  "duplicate command alias: x",
		},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cmds, err := handler.Cmds(tc.cmds)
			assert.Nil(t, cmds)
			assert.EqualError(t, err, tc.err)
		})
	}

	// The same name is fine below different commands.
	cmds, err := handler.Cmds([]*handler.Command{
		{Name: "a", Action: ok, Subcmds: []*handler.Command{{Name: "list", Action: ok}}},
		{Name: "b", Action: ok, Subcmds: []*handler.Command{{Name: "list", Action: ok}}},
	})
	assert.NoError(t, err)
	c, _ := cmds.Get([]string{"b", "list"})
	assert.Equal(t, "list", c.Name)
}

func TestNames(t *testing.T) {
	names := newCmds(t).Names()
	sort.Strings(names)
	assert.Equal(t, []string{"", "foo", "test"}, []string(names))

	empty, err := handler.Cmds(nil)
	assert.NoError(t, err)
	assert.Empty(t, empty.Names())
}

func TestGet(t *testing.T) {
	cmds := newCmds(t)
	tt := map[string]struct {
		path []string
		name string
	}{
		"top":                {[]string{"foo"}, "foo"},
		"sub":                {[]string{"foo", "SubA"}, "SubA"},
		"alias":              {[]string{"f"}, "foo"},
		"alias below parent": {[]string{"sb"}, "SubB"},
		"fallback":           {[]string{""}, ""},
		"missing":            {[]string{"nope"}, ""},
		"missing sub":        {[]string{"foo", "nope"}, ""},
		"alias not first":    {[]string{"foo", "sb"}, ""},
		"empty path":         {nil, ""},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c, h := cmds.Get(tc.path)
			if tc.name == "" && name != "fallback" {
				assert.Nil(t, c)
				assert.Nil(t, h)
				return
			}
			assert.Equal(t, tc.name, c.Name)
			assert.NotNil(t, h)
		})
	}
}

func TestSeek(t *testing.T) {
	cmds := newCmds(t)
	tt := map[string]struct {
		path []string
		name string
		used int
	}{
		"full path":          {[]string{"foo", "SubA"}, "SubA", 2},
		"arguments follow":   {[]string{"foo", "SubA", "extra", "args"}, "SubA", 2},
		"unknown sub":        {[]string{"foo", "nope", "arg"}, "foo", 1},
		"alias":              {[]string{"sb", "arg"}, "SubB", 1},
		"top alias":          {[]string{"f", "3"}, "foo", 1},
		"unknown":            {[]string{"nope", "arg"}, "", 0},
		"nothing to look at": {[]string{}, "", 0},
		"nil":                {nil, "", 0},
	}
	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			c, h, n := cmds.Seek(tc.path)
			assert.Equal(t, tc.used, n)
			if tc.name == "" {
				assert.Nil(t, c)
				assert.Nil(t, h)
				return
			}
			assert.Equal(t, tc.name, c.Name)
			assert.NotNil(t, h)
		})
	}
}

func TestWriter(t *testing.T) {
	cmds := newCmds(t)
	assert.Nil(t, cmds.Writer([]string{"nope"}))

	write := func(w *handler.CommandWriter) string {
		buf := bytes.NewBuffer(nil)
		n, err := w.WriteTo(buf)
		assert.NoError(t, err)
		assert.Equal(t, int64(buf.Len()), n)
		return buf.String()
	}

	// The alias goes before the name, and the empty-named command is left out.
	assert.Equal(t, "   f, foo does foo\n   test", write(cmds.Writer(nil)))

	// Commands with no Subcmds have nothing to list, and an alias finds the
	// command.
	assert.Equal(t, "", write(cmds.Writer([]string{"test"})))
	assert.Equal(t, "   SubA\n   sb, SubB", write(cmds.Writer([]string{"f"})))

	w := cmds.Writer(nil)
	w.Recursive = true
	w.Padding = "--"
	assert.Equal(t, "--f, foo does foo\n----SubA\n----sb, SubB\n--test", write(w))
}

func TestCommandHelpers(t *testing.T) {
	parent := &handler.Command{Name: "parent"}
	parent.AddSub(&handler.Command{Name: "b"})
	parent.AddSub(&handler.Command{Name: "a"})
	assert.Len(t, parent.Subcmds, 2)

	sort.Slice(parent.Subcmds, func(i, j int) bool {
		return handler.CmdNameLT(parent.Subcmds[i], parent.Subcmds[j])
	})
	assert.Equal(t, "a", parent.Subcmds[0].Name)
	assert.Equal(t, "b", parent.Subcmds[1].Name)
}

func TestCommandsSwitch(t *testing.T) {
	s, err := newCmds(t).Switch()
	assert.NoError(t, err)
	a, err := s.Handle(21)
	assert.NoError(t, err)
	assert.Equal(t, "21", a)
	a, err = s.Handle("8")
	assert.NoError(t, err)
	assert.Equal(t, 8, a)

	// The commands that take no argument can't be found by type.
	_, err = s.Handle(nil)
	assert.Equal(t, handler.ErrNoHandler, err)

	// Two commands that take the same type can't both be in a Switch.
	cmds, err := handler.Cmds([]*handler.Command{
		{Name: "one", Action: func(i int) int { return i }},
		{Name: "two", Action: func(i int) int { return i * 2 }},
	})
	assert.NoError(t, err)
	s, err = cmds.Switch()
	assert.Nil(t, s)
	assert.EqualError(t, err, "commands one and two both take int")
}

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
