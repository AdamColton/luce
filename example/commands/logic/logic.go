package logic

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

// HandlerObject holds the state of the demo. Its methods are the commands: the
// method named XHandler, which takes a pointer to a request struct, is the
// command x, and XUsage describes it.
type HandlerObject struct {
	People       []string
	EmptyCounter int
	Timeout      int
	*cli.ExitCloseHandler
	cli.Helper
}

// EC returns the ExitClose, so that HandlerObject is a cli.Commander.
func (ho *HandlerObject) EC() *cli.ExitClose {
	return ho.ExitClose
}

// PersonReq is the request of the person command.
type PersonReq struct {
	Name string
	Age  int
}

// PersonResp is the response of the person command.
type PersonResp struct {
	Name string
}

// PersonHandler adds the name to the People.
func (ho *HandlerObject) PersonHandler(p *PersonReq) *PersonResp {
	ho.People = append(ho.People, p.Name)
	return &PersonResp{
		Name: p.Name,
	}
}

// PersonUsage describes the person command.
func (*HandlerObject) PersonUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create a Person record",
		Alias: "p",
	}
}

// EmptyReq is the request of the empty command, which has no fields.
type EmptyReq struct{}

// EmptyHandler counts how many times it was called and returns the count.
func (ho *HandlerObject) EmptyHandler(e *EmptyReq) int {
	ho.EmptyCounter++
	return ho.EmptyCounter
}

// EmptyUsage describes the empty command.
func (*HandlerObject) EmptyUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Demonstrate empty struct handling",
	}
}

// ListReq is the request of the list command, which has no fields.
type ListReq struct{}

// ListHandler returns the People.
func (ho *HandlerObject) ListHandler(e *ListReq) []string {
	return ho.People
}

// ListUsage describes the list command.
func (*HandlerObject) ListUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List all the People",
		Alias: "l",
	}
}

// TimeoutReq is the request of the timeout command, which has no fields.
type TimeoutReq struct{}

// TimeoutHandler returns nothing, so there is no response and the command
// times out, which is what it is there to show.
func (ho *HandlerObject) TimeoutHandler(e *TimeoutReq) {
}

// SetTimeoutReq is the request of the setTimeout command.
type SetTimeoutReq struct {
	Timeout int
}

// SetTimeoutResp is the response of the setTimeout command.
type SetTimeoutResp struct {
	Timeout int
}

// SetTimeoutHandler sets the Timeout and returns it.
func (ho *HandlerObject) SetTimeoutHandler(e *SetTimeoutReq) *SetTimeoutResp {
	ho.Timeout = e.Timeout
	return &SetTimeoutResp{e.Timeout}
}

// Commands builds the commands from the methods of ho, sorted by name, with an
// empty command that shows the help, list as a sub command of person and q as
// the alias of exit. The exit command is left out, and so is its alias, when the
// ExitClose cannot exit.
func (ho *HandlerObject) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(ho)
	cmds.Set("", &handler.Command{
		Name:  "",
		Usage: "",
		Action: func() *cli.HelpResp {
			return &cli.HelpResp{}
		},
	})
	l, _ := cmds.Pop("list")
	cmds.GetVal("person").AddSub(l)
	if exit := cmds.GetVal("exit"); exit != nil {
		exit.Alias = "q"
	}

	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}
