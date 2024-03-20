package logic

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

type HandlerObject struct {
	People       []string
	EmptyCounter int
	Timeout      int
	*cli.ExitCloseHandler
	cli.Helper
}

func (ho *HandlerObject) EC() *cli.ExitClose {
	return ho.ExitClose
}

type PersonReq struct {
	Name string
	Age  int
}

type PersonResp struct {
	Name string
}

func (ho *HandlerObject) PersonHandler(p *PersonReq) *PersonResp {
	ho.People = append(ho.People, p.Name)
	return &PersonResp{
		Name: p.Name,
	}
}

func (*HandlerObject) PersonUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create a Person record",
		Alias: "p",
	}
}

type EmptyReq struct{}

func (ho *HandlerObject) EmptyHandler(e *EmptyReq) int {
	ho.EmptyCounter++
	return ho.EmptyCounter
}

func (*HandlerObject) EmptyUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Demonstrate empty struct handling",
	}
}

type ListReq struct{}

func (ho *HandlerObject) ListHandler(e *ListReq) []string {
	return ho.People
}

func (*HandlerObject) ListUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List all the People",
		Alias: "l",
	}
}

type TimeoutReq struct{}

func (ho *HandlerObject) TimeoutHandler(e *TimeoutReq) {
	// No response causes a timeout
}

type SetTimeoutReq struct {
	Timeout int
}

type SetTimeoutResp struct {
	Timeout int
}

func (ho *HandlerObject) SetTimeoutHandler(e *SetTimeoutReq) *SetTimeoutResp {
	ho.Timeout = e.Timeout
	return &SetTimeoutResp{e.Timeout}
}

func (ho *HandlerObject) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(ho)
	cmds[""] = &handler.Command{
		Name:  "",
		Usage: "",
		Action: func() *cli.HelpResp {
			return &cli.HelpResp{}
		},
	}
	l, _ := cmds.Pop("list")
	cmds["person"].AddSub(l)
	cmds["exit"].Alias = "q"

	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}
