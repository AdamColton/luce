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

func (*HandlerObject) PersonUsage() string {
	return "Create a Person record"
}

type EmptyReq struct{}

func (ho *HandlerObject) EmptyHandler(e *EmptyReq) int {
	ho.EmptyCounter++
	return ho.EmptyCounter
}

func (*HandlerObject) EmptyUsage() string {
	return "Demonstrate empty struct handling"
}

type ListReq struct{}

func (ho *HandlerObject) ListHandler(e *ListReq) []string {
	return ho.People
}

func (*HandlerObject) ListUsage() string {
	return "List all the People"
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
	cmds = append(cmds, handler.Command{
		Name:  "",
		Usage: "",
		Action: func() *cli.HelpResp {
			return &cli.HelpResp{}
		},
	})

	return lerr.Must(handler.Cmds(cmds))
}
