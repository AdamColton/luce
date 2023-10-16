package handler

import (
	"io"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/lstr"
	"github.com/adamcolton/luce/util/luceio"
)

type Commands struct {
	lookup   map[string]int
	cmds     []Command
	handlers []Handler
	names    []string
}

func Must(commands []Command) *Commands {
	cmds, err := Cmds(commands)
	lerr.Panic(err)
	return cmds
}

func Cmds(commands []Command) (*Commands, error) {
	ln := len(commands)
	out := &Commands{
		cmds:     commands,
		lookup:   make(map[string]int, ln),
		handlers: make([]Handler, ln),
		names:    make([]string, ln),
	}
	for i := range commands {
		c := &(commands[i])
		h, err := New(c.Action, c.Name)
		if err != nil {
			return nil, err
		}
		out.lookup[c.Name] = i
		out.handlers[i] = *h
		out.names[i] = c.Name
	}
	return out, nil
}

func (cs *Commands) Switch() (*Switch, error) {
	s := NewSwitch(len(cs.cmds))
	for i := range cs.cmds {
		s.RegisterHandler(&cs.handlers[i])
	}
	return s, nil
}

func (cs *Commands) Names() []string {
	return cs.names
}

func (cs *Commands) Get(name string) (*Command, *Handler) {
	idx, found := cs.lookup[name]
	if !found {
		return nil, nil
	}
	return &cs.cmds[idx], &cs.handlers[idx]
}

var maxStr = iter.Max(lstr.Len)

func (cs *Commands) WriteTo(w io.Writer) (int64, error) {
	// TODO: add Len function to lstr
	// then similar in slice as func[T any] Len(s []T) int {return len(s)}
	// same in channel
	ln := maxStr.Iter(0, slice.NewIter(cs.names)) + 3

	sw := luceio.NewSumWriter(w)
	for i, n := range cs.names {
		if i > 0 {
			sw.WriteString("\n")
		}
		sw.WriteStrings("   ", n)
		if cmd, _ := cs.Get(n); cmd.Usage != "" {
			sw.WriteStrings(strings.Repeat(" ", ln-len(n)), cmd.Usage)
		}
	}

	return sw.Rets()
}
