package handler

type Command struct {
	Name    string
	Usage   string
	Action  any
	Subcmds []Command
	Alias   string
}
