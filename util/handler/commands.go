package handler

import (
	"io"
	"maps"
	"slices"
	"strings"

	"github.com/adamcolton/luce/ds/idx/hierarchy"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/luceio"
)

// hid identifies a command in the hierarchy. 0 is the root, above the top-level
// commands.
type hid int

// Commands holds a tree of Commands and the Handler for each. Every command is
// found by its path, the names from the top-level command down to it.
type Commands struct {
	cmds     map[hid]*Command
	handlers map[hid]*Handler
	h        *hierarchy.Hierarchy[hid, string]
	alias    map[string]hid
}

// Cmds builds Commands from a tree of Command. It returns an error if an Action
// is not a func that Handler allows, if two commands with the same parent have
// the same name, or if an alias is used more than once.
func Cmds(commands []*Command) (*Commands, error) {
	ln := cmdsLen(commands)
	out := &Commands{
		cmds:     make(map[hid]*Command, ln),
		handlers: make(map[hid]*Handler, ln),
		h:        hierarchy.New[hid, string](ln),
		alias:    make(map[string]hid),
	}
	err := out.addCmds(0, commands)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// cmdsLen counts the commands, including all the Subcmds.
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
		cid, found := cs.h.Key(id, c.Name, true)
		if found {
			return lerr.Str("duplicate command name: " + c.Name)
		}
		h, err := New(c.Action)
		if err != nil {
			return err
		}
		if c.Alias != "" {
			if _, used := cs.alias[c.Alias]; used {
				return lerr.Str("duplicate command alias: " + c.Alias)
			}
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

// Names returns the names of the top-level commands, in no particular order.
func (cs *Commands) Names() slice.Slice[string] {
	return cs.h.Children[0].Slice(nil)
}

// Switch creates a Switch that invokes a command by the type of its argument.
// A command that takes no argument can't be found by type, so it is left out.
// It returns an error if two commands take the same type of argument.
func (cs *Commands) Switch() (*Switch, error) {
	s := NewSwitch(len(cs.cmds))
	// The commands go in by ID, so the same commands always give the same error.
	byType := make(map[any]*Command, len(cs.handlers))
	for _, id := range slices.Sorted(maps.Keys(cs.handlers)) {
		h := cs.handlers[id]
		t := h.Type()
		if t == nil {
			continue
		}
		if other, found := byType[t]; found {
			return nil, lerr.Str("commands " + other.Name + " and " + cs.cmds[id].Name + " both take " + t.String())
		}
		byType[t] = cs.cmds[id]
		s.RegisterHandler(h)
	}
	return s, nil
}

// Get returns the Command and Handler at the path, or nil for both if there is
// none. A path of one name that is not a top-level command is looked up as an
// alias.
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
// Seek and the remainder can be processed as arguments. If the first name is
// not a top-level command it is looked up as an alias. If nothing is found the
// Command and Handler are nil and the int is 0.
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

	if i == 0 && len(path) > 0 {
		cid, found = cs.alias[path[0]]
		if found {
			i = 1
		}
	}

	return cs.cmds[cid], cs.handlers[cid], i
}

// CommandWriter writes a list of the commands below a command, made by
// Commands.Writer.
//
// == projects.Code.luce.handler ==
// [ ] CommandWriter add a newline option, and could it use navigator?
type CommandWriter struct {
	// Padding is written before every command. When Recursive, it is added again
	// for each level of Subcmds.
	Padding string
	// Recursive lists the Subcmds of each command below it.
	Recursive bool
	cid       hid
	cmds      *Commands
}

// Writer returns a CommandWriter for the commands below the one at the path,
// which is looked up as Get does. An empty path is the top-level commands. It
// returns nil if there is no command at the path. Padding starts as 3 spaces.
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

var maxNameLen = liter.Max(func(c *Command) int {
	ln := len(c.Name)
	if c.Alias != "" {
		ln += len(c.Alias) + 2
	}
	return ln
})

// WriteTo writes a line for each command, ordered by name. It has the alias, if
// there is one, then the name, then the usage lined up in a column. Commands
// with an empty name are left out, and there is no newline after the last line.
// It fulfills io.WriterTo.
func (cw *CommandWriter) WriteTo(w io.Writer) (n int64, err error) {
	sw := luceio.NewSumWriter(w)
	cw.writeTo(cw.Padding, "", sw)
	return sw.Rets()
}

func (cw *CommandWriter) writeTo(basePadding, sep string, sw *luceio.SumWriter) {
	cmdNames := cw.cmds.h.Children[cw.cid]
	cmds := make(slice.Slice[*Command], 0, cmdNames.Len())
	cmdNames.Each(func(name string, done *bool) {
		id, _ := cw.cmds.h.Key(cw.cid, name, false)
		cmds = append(cmds, cw.cmds.cmds[id])
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
