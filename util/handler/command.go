package handler

type Command struct {
	Name    string
	Usage   string
	Action  any
	Subcmds []*Command
	Alias   string
}

func (c *Command) AddSub(sub *Command) {
	c.Subcmds = append(c.Subcmds, sub)
}

func CmdNameLT(i, j *Command) bool {
	return i.Name < j.Name
}
