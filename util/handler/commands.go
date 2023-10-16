package handler

import (
	"io"
	"strings"

	"github.com/adamcolton/luce/ds/idx/hierarchy"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/luceio"
)

type hid int

type Commands struct {
	cmds     map[hid]*Command
	handlers map[hid]*Handler
	h        *hierarchy.Hierarchy[hid, string]
	byId     map[hid]int
	alias    map[string]hid
}

func Cmds(commands []*Command) (*Commands, error) {
	// todo: get len w/ subcommands
	ln := cmdsLen(commands)
	out := &Commands{
		cmds:     make(map[hid]*Command, ln),
		handlers: make(map[hid]*Handler, ln),
		h:        hierarchy.New[hid, string](ln),
		byId:     make(map[hid]int, ln),
		alias:    make(map[string]hid),
	}
	err := out.addCmds(0, commands)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func cmdsLen(commands []*Command) int {
	out := len(commands)
	for _, c := range commands {
		if len(c.Subcmds) > 0 {
			out += cmdsLen(c.Subcmds)
		}
	}
	return out
}

func (cs *Commands) addCmds(id hid, commands []*Command) error {
	for _, c := range commands {
		cid, _ := cs.h.Key(id, c.Name, true)
		h, err := New(c.Action)
		if err != nil {
			return err
		}
		if c.Alias != "" {
			cs.alias[c.Alias] = cid
		}
		cs.cmds[cid] = c
		cs.handlers[cid] = h
		err = cs.addCmds(cid, c.Subcmds)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cs *Commands) Names() slice.Slice[string] {
	return cs.h.Children[0].Slice()
}
func (cs *Commands) Switch() (*Switch, error) {
	s := NewSwitch(len(cs.cmds))
	for _, h := range cs.handlers {
		s.RegisterHandler(h)
	}
	return s, nil
}

func (cs *Commands) Get(path []string) (*Command, *Handler) {
	cid, found := cs.h.Get(path, false)
	if !found && len(path) == 1 {
		cid, found = cs.alias[path[0]]
	}
	if !found {
		return nil, nil
	}
	return cs.cmds[cid], cs.handlers[cid]
}

// Seek consumes 'path' until no command is found. The int indicates the
// number of strings used. This allows a full line of input to be passed into
// Seek and the remainder can be processed as arguments.
func (cs *Commands) Seek(path []string) (*Command, *Handler, int) {
	var cid, next hid
	var found bool
	i := 0
	for ; i < len(path); i++ {
		next, found = cs.h.Key(cid, path[i], false)
		if found {
			cid = next
		} else {
			break
		}
	}

	if i == 0 {
		cid, found = cs.alias[path[0]]
		if found {
			i = 1
		}
	}

	return cs.cmds[cid], cs.handlers[cid], i
}

// make public add options (padding, newline, recursive)
// could this use navigator?
type CommandWriter struct {
	Padding   string
	Recursive bool
	cid       hid
	cmds      *Commands
}

func (cs *Commands) Writer(path []string) *CommandWriter {
	cid, found := cs.h.Get(path, false)
	if !found && len(path) == 1 {
		cid, found = cs.alias[path[0]]
	}
	if !found {
		return nil
	}
	return &CommandWriter{
		Padding: "   ",
		cid:     cid,
		cmds:    cs,
	}
}

var maxNameLen = iter.Max(func(c *Command) int {
	ln := len(c.Name)
	if c.Alias != "" {
		ln += len(c.Alias) + 2
	}
	return ln
})

func (cw *CommandWriter) WriteTo(w io.Writer) (n int64, err error) {
	sw := luceio.NewSumWriter(w)
	cw.writeTo(cw.Padding, "", sw)
	return sw.Rets()
}

func (cw *CommandWriter) writeTo(basePadding, sep string, sw *luceio.SumWriter) {
	cmdNames := cw.cmds.h.Children[cw.cid]
	cmds := make(slice.Slice[*Command], 0, cmdNames.Len())
	cmdNames.Each(func(name string) bool {
		id, _ := cw.cmds.h.Key(cw.cid, name, false)
		cmds = append(cmds, cw.cmds.cmds[id])
		return false
	})
	cmds.Sort(CmdNameLT)

	ln := maxNameLen.Iter(0, cmds.Iter()) + 3

	for _, c := range cmds {
		if c.Name == "" {
			continue
		}
		sw.WriteString(sep)
		sep = "\n"
		sw.WriteStrings(cw.Padding)
		col0ln := 0
		if c.Alias != "" {
			sw.WriteStrings(c.Alias, ", ")
			col0ln += len(c.Alias) + 2
		}
		sw.WriteString(c.Name)
		col0ln += len(c.Name) + 2
		if c.Usage != "" {
			sw.WriteStrings(strings.Repeat(" ", ln-col0ln))
			sw.WriteString(c.Usage)
		}
		if cw.Recursive {
			id, _ := cw.cmds.h.Key(cw.cid, c.Name, false)
			(&CommandWriter{
				Padding:   cw.Padding + basePadding,
				Recursive: true,
				cid:       id,
				cmds:      cw.cmds,
			}).writeTo(basePadding, sep, sw)
		}
	}
}
