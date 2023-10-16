package handler

// Command is an action with a name. Commands are arranged in a tree with
// Subcmds, and Cmds turns the tree into Commands.
type Command struct {
	// Name of the command. It is unique among the commands with the same parent.
	// A command with an empty Name is not listed by a CommandWriter. The Runner in
	// util/cli runs it when the input matches no other command.
	Name string
	// Usage is a short description of the command.
	Usage string
	// Action is a func that Handler allows. It is called with the argument that
	// follows the command's name.
	Action any
	// Subcmds are the commands below this one, so they are invoked by the name
	// of this command followed by their own.
	Subcmds []*Command
	// Alias is a shorter name for the command, used in place of the whole path
	// to it. Aliases are shared by all the commands, so each can be used once.
	Alias string
}

// AddSub adds a command below this one.
func (c *Command) AddSub(sub *Command) {
	c.Subcmds = append(c.Subcmds, sub)
}

// CmdNameLT is a less func that orders commands by name.
func CmdNameLT(i, j *Command) bool {
	return i.Name < j.Name
}
